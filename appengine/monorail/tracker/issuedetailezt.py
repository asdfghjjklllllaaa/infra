# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Classes that implement the issue detail page and related forms.

Summary of classes:
  IssueDetailEzt: Show one issue in detail w/ all metadata and comments, and
               process additional comments or metadata changes on it.
  FlagSpamForm: Record the user's desire to report the issue as spam.
"""
from __future__ import print_function
from __future__ import division
from __future__ import absolute_import

import httplib
import json
import logging
import time
from third_party import ezt

import settings
from api import converters
from businesslogic import work_env
from features import features_bizobj
from features import send_notifications
from features import hotlist_helpers
from features import hotlist_views
from framework import exceptions
from framework import framework_bizobj
from framework import framework_constants
from framework import framework_helpers
from framework import framework_views
from framework import jsonfeed
from framework import paginate
from framework import permissions
from framework import servlet
from framework import servlet_helpers
from framework import sorting
from framework import sql
from framework import template_helpers
from framework import urls
from framework import xsrf
from proto import user_pb2
from proto import tracker_pb2
from services import features_svc
from services import tracker_fulltext
from tracker import field_helpers
from tracker import issuepeek
from tracker import tracker_bizobj
from tracker import tracker_constants
from tracker import tracker_helpers
from tracker import tracker_views

from google.protobuf import json_format


class IssueDetailEzt(issuepeek.IssuePeek):
  """IssueDetailEzt is a page that shows the details of one issue."""

  _PAGE_TEMPLATE = 'tracker/issue-detail-page.ezt'
  _MISSING_ISSUE_PAGE_TEMPLATE = 'tracker/issue-missing-page.ezt'
  _MAIN_TAB_MODE = issuepeek.IssuePeek.MAIN_TAB_ISSUES
  _ALLOW_VIEWING_DELETED = True

  def __init__(self, request, response, **kwargs):
    super(IssueDetailEzt, self).__init__(request, response, **kwargs)
    self.missing_issue_template = template_helpers.MonorailTemplate(
        self._TEMPLATE_PATH + self._MISSING_ISSUE_PAGE_TEMPLATE)

  def GetTemplate(self, page_data):
    """Return a custom 404 page for skipped issue local IDs."""
    if page_data.get('http_response_code', httplib.OK) == httplib.NOT_FOUND:
      return self.missing_issue_template
    else:
      return servlet.Servlet.GetTemplate(self, page_data)

  def _GetMissingIssuePageData(
      self, mr, issue_deleted=False, issue_missing=False,
      issue_not_specified=False, issue_not_created=False,
      moved_to_project_name=None, moved_to_id=None,
      local_id=None, page_perms=None, delete_form_token=None):
    if not page_perms:
      # Make a default page perms.
      page_perms = self.MakePagePerms(mr, None, granted_perms=None)
      page_perms.CreateIssue = False
    return {
        'issue_tab_mode': 'issueDetail',
        'http_response_code': httplib.NOT_FOUND,
        'issue_deleted': ezt.boolean(issue_deleted),
        'issue_missing': ezt.boolean(issue_missing),
        'issue_not_specified': ezt.boolean(issue_not_specified),
        'issue_not_created': ezt.boolean(issue_not_created),
        'moved_to_project_name': moved_to_project_name,
        'moved_to_id': moved_to_id,
        'local_id': local_id,
        'page_perms': page_perms,
        'delete_form_token': delete_form_token,
     }

  def _MakeIssueView(
      self, mr, issue, users_by_id, config, issue_reporters=None):
    """Create view objects that help display parts of an issue.

    Args:
      mr: commonly used info parsed from the request.
      issue: issue PB for the currently viewed issue.
      users_by_id: dictionary of {user_id: UserView,...}.
      config: ProjectIssueConfig for the project that contains this issue.
      issue_reporters: list of user IDs who have flagged the issue as spam.

    Returns:
      The IssueView for the whole issue.
    """
    with mr.profiler.Phase('getting related issues'):
      open_related, closed_related = (
          tracker_helpers.GetAllowedOpenAndClosedRelatedIssues(
              self.services, mr, issue))
      all_related_iids = list(issue.blocked_on_iids) + list(issue.blocking_iids)
      if issue.merged_into:
        all_related_iids.append(issue.merged_into)
      all_related = self.services.issue.GetIssues(mr.cnxn, all_related_iids)

    with mr.profiler.Phase('making issue view'):
      issue_view = tracker_views.IssueView(
          issue, users_by_id, config,
          open_related=open_related, closed_related=closed_related,
          all_related={rel.issue_id: rel for rel in all_related})

    issue_reporters = issue_reporters or []
    issue_view.flagged_spam = mr.auth.user_id in issue_reporters

    return issue_view

  def GatherPageData(self, mr):
    """Build up a dictionary of data values to use when rendering the page.

    Args:
      mr: commonly used info parsed from the request.

    Returns:
      Dict of values used by EZT for rendering the page.
    """
    with work_env.WorkEnv(mr, self.services) as we:
      config = we.GetProjectConfig(mr.project_id)

      if mr.local_id is None:
        return self._GetMissingIssuePageData(mr, issue_not_specified=True)
      try:
        # Signed in users could edit the issue, so it must be fresh.
        use_cache = not mr.auth.user_id
        issue = we.GetIssueByLocalID(
            mr.project_id, mr.local_id, use_cache=use_cache,
            allow_viewing_deleted=True)
        comments = we.ListIssueComments(issue)
      except exceptions.NoSuchIssueException:
        issue = None

      # Show explanation of skipped issue local IDs or deleted issues.
      if issue is None or issue.deleted:
        missing = mr.local_id <= self.services.issue.GetHighestLocalID(
            mr.cnxn, mr.project_id)
        if missing or (issue and issue.deleted):
          moved_to_ref = self.services.issue.GetCurrentLocationOfMovedIssue(
              mr.cnxn, mr.project_id, mr.local_id)
          moved_to_project_id, moved_to_id = moved_to_ref
          if moved_to_project_id is not None:
            moved_to_project = we.GetProject(moved_to_project_id)
            moved_to_project_name = moved_to_project.project_name
          else:
            moved_to_project_name = None

          if issue:
            granted_perms = tracker_bizobj.GetGrantedPerms(
                issue, mr.auth.effective_ids, config)
          else:
            granted_perms = None
          page_perms = self.MakePagePerms(
              mr, issue,
              permissions.DELETE_ISSUE, permissions.CREATE_ISSUE,
              granted_perms=granted_perms)
          return self._GetMissingIssuePageData(
              mr,
              issue_deleted=ezt.boolean(issue is not None),
              issue_missing=ezt.boolean(issue is None and missing),
              moved_to_project_name=moved_to_project_name,
              moved_to_id=moved_to_id,
              local_id=mr.local_id,
              page_perms=page_perms)
        else:
          # Issue is not "missing," moved, or deleted, it is just non-existent.
          return self._GetMissingIssuePageData(mr, issue_not_created=True)

      if issue.approval_values:
        logging.info(
            'Approval issues cannot be viewed in the legacy UI.')
        url = framework_helpers.FormatAbsoluteURL(
            mr, urls.ISSUE_DETAIL, id=issue.local_id)
        return self.redirect(url, abort=True)

      star_cnxn = sql.MonorailConnection()
      star_promise = framework_helpers.Promise(
          we.IsIssueStarred, issue, cnxn=star_cnxn)
      userprefs = we.GetUserPrefs(mr.auth.user_id)
      code_font = any(pref for pref in userprefs.prefs
                      if pref.name == 'code_font' and pref.value == 'true')

    granted_perms = tracker_bizobj.GetGrantedPerms(
        issue, mr.auth.effective_ids, config)

    page_perms = self.MakePagePerms(
        mr, issue,
        permissions.CREATE_ISSUE,
        permissions.FLAG_SPAM,
        permissions.VERDICT_SPAM,
        permissions.SET_STAR,
        permissions.EDIT_ISSUE,
        permissions.EDIT_ISSUE_SUMMARY,
        permissions.EDIT_ISSUE_STATUS,
        permissions.EDIT_ISSUE_OWNER,
        permissions.EDIT_ISSUE_CC,
        permissions.DELETE_ISSUE,
        permissions.ADD_ISSUE_COMMENT,
        permissions.DELETE_OWN,
        permissions.DELETE_ANY,
        permissions.VIEW_INBOUND_MESSAGES,
        granted_perms=granted_perms)

    issue_spam_promise = None
    issue_spam_hist_promise = None

    if page_perms.FlagSpam:
      issue_spam_cnxn = sql.MonorailConnection()
      issue_spam_promise = framework_helpers.Promise(
          self.services.spam.LookupIssueFlaggers, issue_spam_cnxn,
          issue.issue_id)

    if page_perms.VerdictSpam:
      issue_spam_hist_cnxn = sql.MonorailConnection()
      issue_spam_hist_promise = framework_helpers.Promise(
          self.services.spam.LookupIssueVerdictHistory, issue_spam_hist_cnxn,
          [issue.issue_id])


    users_involved_in_issue = tracker_bizobj.UsersInvolvedInIssues([issue])
    users_involved_in_comment_list = tracker_bizobj.UsersInvolvedInCommentList(
        comments)
    with mr.profiler.Phase('making user views'):
      users_by_id = framework_views.MakeAllUserViews(
          mr.cnxn, self.services.user, users_involved_in_issue,
          users_involved_in_comment_list)
      framework_views.RevealAllEmailsToMembers(mr.auth, mr.project, users_by_id)

    issue_flaggers, comment_flaggers = [], {}
    if issue_spam_promise:
      issue_flaggers, comment_flaggers = issue_spam_promise.WaitAndGetValue()

    issue_view = self._MakeIssueView(
        mr, issue, users_by_id, config, issue_flaggers)

    with mr.profiler.Phase('converting comments to ListCommentsResponse'):
      issue_perms = permissions.UpdateIssuePermissions(
          mr.perms, mr.project, issue, mr.auth.effective_ids, config=config)
      comments_list = converters.ConvertCommentList(
          issue, comments, config, users_by_id, comment_flaggers,
          mr.auth.user_id, issue_perms)
      comments_list = [
          json_format.MessageToDict(comment) for comment in comments_list]

    with mr.profiler.Phase('getting starring info'):
      starred = star_promise.WaitAndGetValue()
      star_cnxn.Close()
      permit_edit = permissions.CanEditIssue(
          mr.auth.effective_ids, mr.perms, mr.project, issue,
          granted_perms=granted_perms)
      page_perms.EditIssue = ezt.boolean(permit_edit)
      permit_edit_cc = self.CheckPerm(
          mr, permissions.EDIT_ISSUE_CC, art=issue, granted_perms=granted_perms)
      discourage_plus_one = not (starred or permit_edit or permit_edit_cc)

    # Check whether to allow attachments from the details page
    allow_attachments = tracker_helpers.IsUnderSoftAttachmentQuota(mr.project)

    hotlist_id = mr.GetIntParam('hotlist_id', None)
    hotlist = None
    if hotlist_id:
      try:
        hotlist = self.services.features.GetHotlist(mr.cnxn, hotlist_id)
      except features_svc.NoSuchHotlistException:
        pass

    if hotlist:
      mr.ComputeColSpec(hotlist)
    else:
      mr.ComputeColSpec(config)

    restrict_to_known = config.restrict_to_known
    field_name_set = {fd.field_name.lower() for fd in config.field_defs
                      if fd.field_type is tracker_pb2.FieldTypes.ENUM_TYPE and
                      not fd.is_deleted}  # TODO(jrobbins): restrictions
    non_masked_labels = tracker_bizobj.NonMaskedLabels(
        issue.labels, field_name_set)

    component_paths = []
    for comp_id in issue.component_ids:
      cd = tracker_bizobj.FindComponentDefByID(comp_id, config)
      if cd:
        component_paths.append(cd.path)
      else:
        logging.warn(
            'Issue %r has unknown component %r', issue.issue_id, comp_id)
    initial_components = ', '.join(component_paths)

    after_issue_update = tracker_constants.DEFAULT_AFTER_ISSUE_UPDATE
    if mr.auth.user_pb:
      after_issue_update = mr.auth.user_pb.after_issue_update

    prevent_restriction_removal = (
        mr.project.only_owners_remove_restrictions and
        not framework_bizobj.UserOwnsProject(
            mr.project, mr.auth.effective_ids))

    offer_issue_copy_move = True
    for lab in tracker_bizobj.GetLabels(issue):
      if lab.lower().startswith('restrict-'):
        offer_issue_copy_move = False

    previous_locations = self.GetPreviousLocations(mr, issue)

    spam_verdict_history = []
    if issue_spam_hist_promise:
      issue_spam_hist = issue_spam_hist_promise.WaitAndGetValue()

      spam_verdict_history = [template_helpers.EZTItem(
          created=verdict['created'].isoformat(),
          is_spam=verdict['is_spam'],
          reason=verdict['reason'],
          user_id=verdict['user_id'],
          classifier_confidence=verdict['classifier_confidence'],
          overruled=verdict['overruled'],
          ) for verdict in issue_spam_hist]

    # get hotlists that contain the current issue
    issue_hotlists = self.services.features.GetHotlistsByIssueID(
        mr.cnxn, issue.issue_id)
    users_by_id = framework_views.MakeAllUserViews(
        mr.cnxn, self.services.user, features_bizobj.UsersInvolvedInHotlists(
            issue_hotlists))

    issue_hotlist_views = [hotlist_views.HotlistView(
        hotlist_pb, mr.perms, mr.auth, mr.auth.user_id, users_by_id,
        self.services.hotlist_star.IsItemStarredBy(
            mr.cnxn, hotlist_pb.hotlist_id, mr.auth.user_id)
    ) for hotlist_pb in self.services.features.GetHotlistsByIssueID(
        mr.cnxn, issue.issue_id)]

    visible_issue_hotlist_views = [view for view in issue_hotlist_views if
                                   view.visible]

    (user_issue_hotlist_views, involved_users_issue_hotlist_views,
     remaining_issue_hotlist_views) = _GetBinnedHotlistViews(
         visible_issue_hotlist_views, users_involved_in_issue)

    user_remaining_hotlists = [hotlist for hotlist in
                      self.services.features.GetHotlistsByUserID(
                          mr.cnxn, mr.auth.user_id) if
                      hotlist not in issue_hotlists]

    is_member = framework_bizobj.UserIsInProject(
        mr.project, mr.auth.effective_ids)

    description_list = [
        comment for i, comment in enumerate(comments_list)
        if 'descriptionNum' in comment or i == 0]

    reporter_name = comments_list[0]['commenter']['displayName']
    reporter_user_id = comments_list[0]['commenter'].get('userId', 0)
    reported_timestamp = comments_list[0]['timestamp']

    return {
        'comment_list': json.dumps(comments_list[1:]),
        'description_list': json.dumps(description_list),
        'reporter_name': reporter_name,
        'reporter_user_id': reporter_user_id,
        'reported_timestamp': reported_timestamp,
        'issue_tab_mode': 'issueDetail',
        'issue': issue_view,
        'title_summary': issue_view.summary,  # used in <head><title>
        'noisy': ezt.boolean(tracker_helpers.IsNoisy(
            len(comments_list), issue.star_count)),
        'link_rel_canonical': framework_helpers.FormatCanonicalURL(mr, ['id']),

        'flipper_hotlist_id': hotlist_id,
        'searchtip': 'You can jump to any issue by number',
        'starred': ezt.boolean(starred),
        'discourage_plus_one': ezt.boolean(discourage_plus_one),
        'pagegen': str(int(time.time() * 1000000)),

        # For deep linking and input correction after a failed submit.
        'initial_summary': issue_view.summary,
        'initial_comment': '',
        'initial_status': issue_view.status.name,
        'initial_owner': issue_view.owner.email,
        'initial_cc': ', '.join([pb.email for pb in issue_view.cc]),
        'initial_blocked_on': issue_view.blocked_on_str,
        'initial_blocking': issue_view.blocking_str,
        'initial_merge_into': issue_view.merged_into_str,
        'labels': non_masked_labels,
        'initial_components': initial_components,
        'fields': issue_view.fields,

        'any_errors': ezt.boolean(mr.errors.AnyErrors()),
        'allow_attachments': ezt.boolean(allow_attachments),
        'max_attach_size': template_helpers.BytesKbOrMb(
            framework_constants.MAX_POST_BODY_SIZE),
        'colspec': mr.col_spec,
        'restrict_to_known': ezt.boolean(restrict_to_known),
        'after_issue_update': int(after_issue_update),  # TODO(jrobbins): str
        'prevent_restriction_removal': ezt.boolean(
            prevent_restriction_removal),
        'offer_issue_copy_move': ezt.boolean(offer_issue_copy_move),
        'statuses_offer_merge': config.statuses_offer_merge,
        'page_perms': page_perms,
        'previous_locations': previous_locations,
        'spam_verdict_history': spam_verdict_history,

        # For adding issue to user's hotlists
        'user_remaining_hotlists': user_remaining_hotlists,
        # For showing hotlists that contain this issue
        'user_issue_hotlists': user_issue_hotlist_views,
        'involved_users_issue_hotlists': involved_users_issue_hotlist_views,
        'remaining_issue_hotlists': remaining_issue_hotlist_views,

        'is_member': ezt.boolean(is_member),
        'code_font': ezt.boolean(code_font),

        'other_ui_path': 'issues/detail',
        'local_id': issue_view.local_id,
    }

  def GatherHelpData(self, mr, page_data):
    """Return a dict of values to drive on-page user help.

    Args:
      mr: commonly used info parsed from the request.
      page_data: Dictionary of base and page template data.

    Returns:
      A dict of values to drive on-page user help, to be added to page_data.
    """
    help_data = super(IssueDetailEzt, self).GatherHelpData(mr, page_data)
    dismissed = []
    if mr.auth.user_pb:
      with work_env.WorkEnv(mr, self.services) as we:
        userprefs = we.GetUserPrefs(mr.auth.user_id)
      dismissed = [
          pv.name for pv in userprefs.prefs if pv.value == 'true']

    is_privileged_domain_user = framework_bizobj.IsPriviledgedDomainUser(
        mr.auth.user_pb.email)
    # Check if the user's query is just the ID of an existing issue.
    # If so, display a "did you mean to search?" cue card.
    jump_local_id = None
    any_availability_message = False
    iv = page_data.get('issue')
    if iv:
      participant_views = (
          [iv.owner, iv.derived_owner] + iv.cc + iv.derived_cc)
      any_availability_message = any(
          pv.avail_message for pv in participant_views
          if pv and pv.user_id)

    if (mr.auth.user_id and
        'privacy_click_through' not in dismissed):
      help_data['cue'] = 'privacy_click_through'
    elif (mr.auth.user_id and
        'code_of_conduct' not in dismissed):
      help_data['cue'] = 'code_of_conduct'
    elif (tracker_constants.JUMP_RE.match(mr.query) and
          'search_for_numbers' not in dismissed):
      jump_local_id = int(mr.query)
      help_data['cue'] = 'search_for_numbers'
    elif (any_availability_message and
          'availability_msgs' not in dismissed):
      help_data['cue'] = 'availability_msgs'

    help_data.update({
        'is_privileged_domain_user': ezt.boolean(is_privileged_domain_user),
        'jump_local_id': jump_local_id,
        })
    return help_data

  # TODO(sheyang): Support comments incremental loading in API
  def _PaginatePartialComments(self, mr, issue):
    """Load and paginate the visible comments for the given issue."""
    abbr_comment_rows = self.services.issue.GetAbbrCommentsForIssue(
          mr.cnxn, issue.issue_id)
    if not abbr_comment_rows:
      return [], [], None

    comments = abbr_comment_rows[1:]
    all_comment_ids = [row[0] for row in comments]

    pagination_url = '%s?id=%d' % (urls.ISSUE_DETAIL, issue.local_id)
    url_params = [(name, mr.GetParam(name)) for name in
                 framework_helpers.RECOGNIZED_PARAMS]
    pagination = paginate.VirtualPagination(
        len(all_comment_ids),
        mr.GetPositiveIntParam(
            'cnum', framework_constants.DEFAULT_COMMENTS_PER_PAGE),
        mr.GetPositiveIntParam('cstart'),
        list_page_url=pagination_url, project_name=mr.project_name,
        count_up=False, start_param_name='cstart', num_param_name='cnum',
        max_num=settings.max_comments_per_page, url_params=url_params)
    if pagination.last == 1 and pagination.start == len(all_comment_ids):
      pagination.visible = ezt.boolean(False)

    visible_comment_ids = all_comment_ids[pagination.last - 1:pagination.start]
    visible_comment_seqs = list(range(pagination.last, pagination.start + 1))
    visible_comments = self.services.issue.GetCommentsByID(
          mr.cnxn, visible_comment_ids, visible_comment_seqs)

    # TODO(lukasperaza): update first comments to is_description=TRUE
    # so [abbr_comment_rows[0][0]] can be removed
    description_ids = list(set([abbr_comment_rows[0][0]] +
                           [row[0] for row in abbr_comment_rows if row[3]]))
    description_seqs = [0]
    for i, abbr_comment in enumerate(comments):
      if abbr_comment[3]:
        description_seqs.append(i + 1)
    # TODO(lukasperaza): only get descriptions which haven't been deleted
    descriptions = self.services.issue.GetCommentsByID(
          mr.cnxn, description_ids, description_seqs)

    for i, desc in enumerate(descriptions):
      desc.description_num = str(i + 1)

    return descriptions, visible_comments, pagination


  def _ValidateOwner(self, mr, post_data_owner, parsed_owner_id,
                     original_issue_owner_id):
    """Validates that the issue's owner was changed and is a valid owner.

    Args:
      mr: Commonly used info parsed from the request.
      post_data_owner: The owner as specified in the request's data.
      parsed_owner_id: The owner_id from the request.
      original_issue_owner_id: The original owner id of the issue.

    Returns:
      String error message if the owner fails validation else returns None.
    """
    parsed_owner_valid, msg = tracker_helpers.IsValidIssueOwner(
        mr.cnxn, mr.project, parsed_owner_id, self.services)
    if not parsed_owner_valid:
      # Only fail validation if the user actually changed the email address.
      original_issue_owner = self.services.user.LookupUserEmail(
          mr.cnxn, original_issue_owner_id)
      if post_data_owner != original_issue_owner:
        return msg
      else:
        # The user did not change the owner, thus do not fail validation.
        # See https://bugs.chromium.org/p/monorail/issues/detail?id=28 for
        # more details.
        pass

  def _ValidateCC(self, cc_ids, cc_usernames):
    """Validate cc list."""
    if None in cc_ids:
      invalid_cc = [cc_name for cc_name, cc_id in zip(cc_usernames, cc_ids)
                    if cc_id is None]
      return 'Invalid Cc username: %s' % ', '.join(invalid_cc)

  def ProcessFormData(self, mr, post_data):
    """Process the posted issue update form.

    Args:
      mr: commonly used info parsed from the request.
      post_data: The post_data dict for the current request.

    Returns:
      String URL to redirect the user to after processing.
    """
    with work_env.WorkEnv(mr, self.services) as we:
      issue = we.GetIssueByLocalID(
          mr.project_id, mr.local_id, use_cache=False)

    # Check that the user is logged in; anon users cannot update issues.
    if not mr.auth.user_id:
      logging.info('user was not logged in, cannot update issue')
      raise permissions.PermissionException(
          'User must be logged in to update an issue')

    # Check that the user has permission to add a comment, and to enter
    # metadata if they are trying to do that.
    if not self.CheckPerm(mr, permissions.ADD_ISSUE_COMMENT,
                          art=issue):
      logging.info('user has no permission to add issue comment')
      raise permissions.PermissionException(
          'User has no permission to comment on issue')

    parsed = tracker_helpers.ParseIssueRequest(
        mr.cnxn, post_data, self.services, mr.errors, issue.project_name)
    config = self.services.config.GetProjectConfig(mr.cnxn, mr.project_id)
    bounce_labels = parsed.labels[:]
    bounce_fields = tracker_views.MakeBounceFieldValueViews(
        parsed.fields.vals, parsed.fields.phase_vals, config)
    field_helpers.ShiftEnumFieldsIntoLabels(
        parsed.labels, parsed.labels_remove,
        parsed.fields.vals, parsed.fields.vals_remove, config)
    field_values = field_helpers.ParseFieldValues(
        mr.cnxn, self.services.user, parsed.fields.vals,
        parsed.fields.phase_vals, config)

    component_ids = tracker_helpers.LookupComponentIDs(
        parsed.components.paths, config, mr.errors)

    granted_perms = tracker_bizobj.GetGrantedPerms(
        issue, mr.auth.effective_ids, config)
    # We process edits iff the user has permission, and the form
    # was generated including the editing fields.
    permit_edit = (
        permissions.CanEditIssue(
            mr.auth.effective_ids, mr.perms, mr.project, issue,
            granted_perms=granted_perms) and
        'fields_not_offered' not in post_data)
    page_perms = self.MakePagePerms(
        mr, issue,
        permissions.CREATE_ISSUE,
        permissions.EDIT_ISSUE_SUMMARY,
        permissions.EDIT_ISSUE_STATUS,
        permissions.EDIT_ISSUE_OWNER,
        permissions.EDIT_ISSUE_CC,
        granted_perms=granted_perms)
    page_perms.EditIssue = ezt.boolean(permit_edit)

    if not permit_edit:
      if not _FieldEditPermitted(
          parsed.labels, parsed.blocked_on.entered_str,
          parsed.blocking.entered_str, parsed.summary,
          parsed.status, parsed.users.owner_id,
          parsed.users.cc_ids, page_perms):
        raise permissions.PermissionException(
            'User lacks permission to edit fields')

    page_generation_time = int(post_data['pagegen'])
    reporter_id = mr.auth.user_id
    # TODO(jrobbins): consider captcha 3 score in API

    error_msg = self._ValidateOwner(
        mr, post_data.get('owner', '').strip(), parsed.users.owner_id,
        issue.owner_id)
    if error_msg:
      mr.errors.owner = error_msg

    error_msg = self._ValidateCC(
        parsed.users.cc_ids, parsed.users.cc_usernames)
    if error_msg:
      mr.errors.cc = error_msg

    if len(parsed.comment) > tracker_constants.MAX_COMMENT_CHARS:
      mr.errors.comment = 'Comment is too long'
    logging.info('parsed.summary is %r', parsed.summary)
    if len(parsed.summary) > tracker_constants.MAX_SUMMARY_CHARS:
      mr.errors.summary = 'Summary is too long'

    old_owner_id = tracker_bizobj.GetOwnerId(issue)

    orig_merged_into_iid = issue.merged_into
    merge_into_iid = issue.merged_into
    merge_into_text, merge_into_issue = tracker_helpers.ParseMergeFields(
        mr.cnxn, self.services, mr.project_name, post_data,
        parsed.status, config, issue, mr.errors)
    if merge_into_issue:
      merge_into_iid = merge_into_issue.issue_id
      merge_into_project = self.services.project.GetProjectByName(
          mr.cnxn, merge_into_issue.project_name)
      merge_allowed = tracker_helpers.IsMergeAllowed(
          merge_into_issue, mr, self.services)

      new_starrers = tracker_helpers.GetNewIssueStarrers(
          mr.cnxn, self.services, issue.issue_id, merge_into_iid)

    # For any fields that the user does not have permission to edit, use
    # the current values in the issue rather than whatever strings were parsed.
    labels = parsed.labels
    summary = parsed.summary
    is_description = parsed.is_description
    status = parsed.status
    owner_id = parsed.users.owner_id
    cc_ids = parsed.users.cc_ids
    blocked_on_iids = [iid for iid in parsed.blocked_on.iids
                       if iid != issue.issue_id]
    blocking_iids = [iid for iid in parsed.blocking.iids
                     if iid != issue.issue_id]
    dangling_blocked_on_refs = [tracker_bizobj.MakeDanglingIssueRef(*ref)
                                for ref in parsed.blocked_on.dangling_refs]
    dangling_blocking_refs = [tracker_bizobj.MakeDanglingIssueRef(*ref)
                              for ref in parsed.blocking.dangling_refs]
    if not permit_edit:
      is_description = False
      labels = issue.labels
      field_values = issue.field_values
      component_ids = issue.component_ids
      blocked_on_iids = issue.blocked_on_iids
      blocking_iids = issue.blocking_iids
      dangling_blocked_on_refs = issue.dangling_blocked_on_refs
      dangling_blocking_refs = issue.dangling_blocking_refs
      merge_into_iid = issue.merged_into
      if not page_perms.EditIssueSummary:
        summary = issue.summary
      if not page_perms.EditIssueStatus:
        status = issue.status
      if not page_perms.EditIssueOwner:
        owner_id = issue.owner_id
      if not page_perms.EditIssueCc:
        cc_ids = issue.cc_ids

    field_helpers.ValidateCustomFields(
        mr, self.services, field_values, config, mr.errors)

    hotlist_id = post_data.get('hotlist_id')
    if hotlist_id:
      hotlist_id = int(hotlist_id)

    # Generate redirect URLs before applying issue changes.
    with work_env.WorkEnv(mr, self.services) as we:
      redirect_issue = _GetRedirectIssue(we, issue, hotlist_id)

    orig_blocked_on = issue.blocked_on_iids
    if not mr.errors.AnyErrors():
      with work_env.WorkEnv(mr, self.services) as we:
        try:
          if parsed.attachments:
            new_bytes_used = tracker_helpers.ComputeNewQuotaBytesUsed(
                mr.project, parsed.attachments)
            # TODO(jrobbins): Make quota be calculated and stored as
            # part of applying the comment.
            self.services.project.UpdateProject(
                mr.cnxn, mr.project.project_id,
                attachment_bytes_used=new_bytes_used)

          # Store everything we got from the form.  If the user lacked perms
          # any attempted edit would be a no-op because of the logic above.
          amendments, comment = self.services.issue.ApplyIssueComment(
            mr.cnxn, self.services,
            mr.auth.user_id, mr.project_id, mr.local_id, summary, status,
            owner_id, cc_ids, labels, field_values, component_ids,
            blocked_on_iids, blocking_iids, dangling_blocked_on_refs,
            dangling_blocking_refs, merge_into_iid, index_now=False,
            page_gen_ts=page_generation_time, comment=parsed.comment,
            is_description=is_description, attachments=parsed.attachments,
            kept_attachments=parsed.kept_attachments if is_description else [])
          self.services.project.UpdateRecentActivity(
              mr.cnxn, mr.project.project_id)

          # Also update the Issue PB we have in RAM so that the correct
          # CC list will be used for an issue merge.
          # TODO(jrobbins): refactor the call above to: 1. compute the updates
          # and update the issue PB in RAM, then 2. store the updated issue.
          issue.cc_ids = cc_ids
          issue.labels = labels

        except tracker_helpers.OverAttachmentQuota:
          mr.errors.attachments = 'Project attachment quota exceeded.'

      if (merge_into_issue and merge_into_iid != orig_merged_into_iid and
          merge_allowed):
        tracker_helpers.AddIssueStarrers(
            mr.cnxn, self.services, mr,
            merge_into_iid, merge_into_project, new_starrers)
        merge_comment_pb = tracker_helpers.MergeCCsAndAddComment(
            self.services, mr, issue, merge_into_issue)
      elif merge_into_issue:
        merge_comment_pb = None
        logging.info('merge denied: target issue %s not modified',
                     merge_into_iid)
      # TODO(jrobbins): distinguish between EditIssue and
      # AddIssueComment and do just the part that is allowed.
      # And, give feedback in the source issue if any part of the
      # merge was not allowed.  Maybe use AJAX to check as the
      # user types in the issue local ID.

    copy_to_project = CheckCopyIssueRequest(
        self.services, mr, issue, post_data.get('more_actions') == 'copy',
        post_data.get('copy_to'), mr.errors)
    move_to_project = CheckMoveIssueRequest(
        self.services, mr, issue, post_data.get('more_actions') == 'move',
        post_data.get('move_to'), mr.errors)

    if mr.errors.AnyErrors():
      self.PleaseCorrect(
          mr, initial_summary=parsed.summary,
          initial_status=parsed.status,
          initial_owner=parsed.users.owner_username,
          initial_cc=', '.join(parsed.users.cc_usernames),
          initial_components=', '.join(parsed.components.paths),
          initial_comment=parsed.comment,
          labels=bounce_labels, fields=bounce_fields,
          initial_blocked_on=parsed.blocked_on.entered_str,
          initial_blocking=parsed.blocking.entered_str,
          initial_merge_into=merge_into_text)
      return

    send_email = 'send_email' in post_data or not permit_edit

    moved_to_project_name_and_local_id = None
    copied_to_project_name_and_local_id = None
    if move_to_project:
      moved_to_project_name_and_local_id = self.HandleCopyOrMove(
          mr.cnxn, mr, move_to_project, issue, send_email, move=True)
    elif copy_to_project:
      copied_to_project_name_and_local_id = self.HandleCopyOrMove(
          mr.cnxn, mr, copy_to_project, issue, send_email, move=False)

    if amendments or parsed.comment.strip() or parsed.attachments:
      # TODO(jrobbins): Remove the seq_num parameter after we have
      # deployed the change that switches to comment_id.
      send_notifications.PrepareAndSendIssueChangeNotification(
          issue.issue_id, mr.request.host, reporter_id,
          send_email=send_email, old_owner_id=old_owner_id,
          comment_id=comment.id)

    if merge_into_issue and merge_allowed and merge_comment_pb:
      # TODO(jrobbins): Remove the merge_seq_num parameter after we have
      # deployed the change that switches to comment_id.
      send_notifications.PrepareAndSendIssueChangeNotification(
          merge_into_issue.issue_id, mr.request.host, reporter_id,
          send_email=send_email, comment_id=merge_comment_pb.id)

    if permit_edit:
      # Only users who can edit metadata could have edited blocking.
      blockers_added, blockers_removed = framework_helpers.ComputeListDeltas(
          orig_blocked_on, blocked_on_iids)
      delta_blockers = blockers_added + blockers_removed
      send_notifications.PrepareAndSendIssueBlockingNotification(
          issue.issue_id, mr.request.host,
          delta_blockers, reporter_id, send_email=send_email)
      # We don't send notification emails to newly blocked issues: either they
      # know they are blocked, or they don't care and can be fixed anyway.
      # This is the same behavior as the issue entry page.

    after_issue_update = _DetermineAndSetAfterIssueUpdate(
        self.services, mr, post_data)
    return _Redirect(
        mr, post_data, issue.local_id, config, redirect_issue, hotlist_id,
        moved_to_project_name_and_local_id,
        copied_to_project_name_and_local_id, after_issue_update)

  def HandleCopyOrMove(self, cnxn, mr, dest_project, issue, send_email, move):
    """Handle Requests dealing with copying or moving an issue between projects.

    Args:
      cnxn: connection to the database.
      mr: commonly used info parsed from the request.
      dest_project: The project protobuf we are moving the issue to.
      issue: The issue protobuf being moved.
      send_email: True to send email for these actions.
      move: Whether this is a move request. The original issue will not exist if
            this is True.

    Returns:
      A tuple of (project_id, local_id) of the newly copied / moved issue.
    """
    old_text_ref = 'issue %s:%s' % (issue.project_name, issue.local_id)
    if move:
      tracker_fulltext.UnindexIssues([issue.issue_id])
      moved_back_iids = self.services.issue.MoveIssues(
          cnxn, dest_project, [issue], self.services.user)
      ret_project_name_and_local_id = (issue.project_name, issue.local_id)
      new_text_ref = 'issue %s:%s' % ret_project_name_and_local_id
      if issue.issue_id in moved_back_iids:
        content = 'Moved %s back to %s again.' % (old_text_ref, new_text_ref)
      else:
        content = 'Moved %s to now be %s.' % (old_text_ref, new_text_ref)
      comment = self.services.issue.CreateIssueComment(
          mr.cnxn, issue, mr.auth.user_id, content, amendments=[
              tracker_bizobj.MakeProjectAmendment(dest_project.project_name)])
    else:
      copied_issues = self.services.issue.CopyIssues(
          cnxn, dest_project, [issue], self.services.user, mr.auth.user_id)
      copied_issue = copied_issues[0]
      ret_project_name_and_local_id = (copied_issue.project_name,
                                       copied_issue.local_id)
      new_text_ref = 'issue %s:%s' % ret_project_name_and_local_id

      # Add comment to the copied issue.
      old_issue_content = 'Copied %s to %s' % (old_text_ref, new_text_ref)
      self.services.issue.CreateIssueComment(
          mr.cnxn, issue, mr.auth.user_id, old_issue_content)

      # Add comment to the newly created issue.
      # Add project amendment only if the project changed.
      amendments = []
      if issue.project_id != copied_issue.project_id:
        amendments.append(
            tracker_bizobj.MakeProjectAmendment(dest_project.project_name))
      new_issue_content = 'Copied %s from %s' % (new_text_ref, old_text_ref)
      comment = self.services.issue.CreateIssueComment(
          mr.cnxn, copied_issue,
          mr.auth.user_id, new_issue_content, amendments=amendments)

    tracker_fulltext.IndexIssues(
        mr.cnxn, [issue], self.services.user, self.services.issue,
        self.services.config)

    if send_email:
      logging.info('TODO(jrobbins): send email for a move? or combine? %r',
                   comment)

    return ret_project_name_and_local_id


def _DetermineAndSetAfterIssueUpdate(services, mr, post_data):
  after_issue_update = tracker_constants.DEFAULT_AFTER_ISSUE_UPDATE
  if 'after_issue_update' in post_data:
    after_issue_update = user_pb2.IssueUpdateNav(
        int(post_data['after_issue_update'][0]))
    if after_issue_update != mr.auth.user_pb.after_issue_update:
      logging.info('setting after_issue_update to %r', after_issue_update)
      services.user.UpdateUserSettings(
          mr.cnxn, mr.auth.user_id, mr.auth.user_pb,
          after_issue_update=after_issue_update)

  return after_issue_update


def _Redirect(
    mr, post_data, local_id, config, redirect_issue,
    hotlist_id, moved_to_project_name_and_local_id,
    copied_to_project_name_and_local_id, after_issue_update):
  """Prepare a redirect URL for the issuedetail servlets.

  Args:
    mr: common information parsed from the HTTP request.
    post_data: The post_data dict for the current request.
    local_id: int Issue ID for the current request.
    config: The ProjectIssueConfig pb for the current request.
    redirect_issue: The next issue to redirect to, if any.
    hotlist_id: The current hotlist, if any.
    moved_to_project_name_and_local_id: tuple containing the project name the
      issue was moved to and the local id in that project.
    copied_to_project_name_and_local_id: tuple containing the project name the
      issue was copied to and the local id in that project.
    after_issue_update: User preference on where to go next.

  Returns:
    String URL to redirect the user to after processing.
  """
  mr.can = int(post_data['can'])
  mr.query = post_data['q']
  mr.col_spec = post_data['colspec']
  mr.sort_spec = post_data['sort']
  mr.group_by_spec = post_data['groupby']
  mr.start = int(post_data['start'])
  mr.num = int(post_data['num'])
  mr.local_id = local_id

  if redirect_issue:
    next_id = redirect_issue.local_id
    next_project = redirect_issue.project_name
  else:
    next_id = None
    next_project = None

  # Format a redirect url.
  url = _ChooseNextPage(
      mr, local_id, config, moved_to_project_name_and_local_id,
      copied_to_project_name_and_local_id, after_issue_update, next_id,
      next_project=next_project, hotlist_id=hotlist_id)
  logging.debug('Redirecting user to: %s', url)
  return url


def _GetRedirectIssue(we, current_issue, hotlist_id=None):
  """Get the next issue for current_issue.

  Args:
    we: A WorkEnv instance.
    current_issue: The issue from which to look.
    hotlsit_id (optional): The current hotlist.

  Returns:
    The next issue if found, else None.
  """
  hotlist = None
  if hotlist_id:
    try:
      hotlist = we.services.features.GetHotlist(we.mc.cnxn, hotlist_id)
    except features_svc.NoSuchHotlistException:
      pass

  next_issue = None
  try:
    next_issue = GetAdjacentIssue(we, current_issue,
        hotlist=hotlist, next_issue=True)
  except exceptions.NoSuchIssueException:
    pass

  return next_issue


def _FieldEditPermitted(
    labels, blocked_on_str, blocking_str, summary, status, owner_id, cc_ids,
    page_perms):
  """Check permissions on editing individual form fields.

  This check is only done if the user does not have the overall
  EditIssue perm.  If the user edited any field that they do not have
  permission to edit, then they could have forged a post, or maybe
  they had a valid form open in a browser tab while at the same time
  their perms in the project were reduced.  Either way, the servlet
  gives them a BadRequest HTTP error and makes them go back and try
  again.

  TODO(jrobbins): It would be better to show a custom error page that
  takes the user back to the issue with a new page load rather than
  having the user use the back button.

  Args:
    labels: list of label values parsed from the form.
    blocked_on_str: list of blocked-on values parsed from the form.
    blocking_str: list of blocking values parsed from the form.
    summary: issue summary string parsed from the form.
    status: issue status string parsed from the form.
    owner_id: issue owner user ID parsed from the form and looked up.
    cc_ids: list of user IDs for Cc'd users parsed from the form.
    page_perms: object with fields for permissions the current user
        has on the current issue.

  Returns:
    True if there was no permission violation.  False if the user tried
    to edit something that they do not have permission to edit.
  """
  if labels or blocked_on_str or blocking_str:
    logging.info('user has no permission to edit issue metadata')
    return False

  if summary and not page_perms.EditIssueSummary:
    logging.info('user has no permission to edit issue summary field')
    return False

  if status and not page_perms.EditIssueStatus:
    logging.info('user has no permission to edit issue status field')
    return False

  if owner_id and not page_perms.EditIssueOwner:
    logging.info('user has no permission to edit issue owner field')
    return False

  if cc_ids and not page_perms.EditIssueCc:
    logging.info('user has no permission to edit issue cc field')
    return False

  return True


def _ChooseNextPage(
    mr, local_id, config, moved_to_project_name_and_local_id,
    copied_to_project_name_and_local_id, after_issue_update, next_id,
    next_project=None, hotlist_id=None):
  """Choose the next page to show the user after an issue update.

  Args:
    mr: information parsed from the request.
    local_id: int Issue ID of the issue that was updated.
    config: project issue config object.
    moved_to_project_name_and_local_id: tuple containing the project name the
      issue was moved to and the local id in that project.
    copied_to_project_name_and_local_id: tuple containing the project name the
      issue was copied to and the local id in that project.
    after_issue_update: user pref on where to go next.
    next_id: string local ID of next issue at the time the form was generated.
    next_project: project name of the next issue's project, None if next
      issue's project is the same as the current's project (before any changes)
    hotlist_id: optional hotlist_id for when an issue is visited via a hotlist

  Returns:
    String absolute URL of next page to view.
  """
  issue_ref_str = '%s:%d' % (mr.project_name, local_id)
  if next_project is None:
    next_project = mr.project_name
  kwargs = {
    'ts': int(time.time()),
    'cursor': issue_ref_str,
  }
  if moved_to_project_name_and_local_id:
    kwargs['moved_to_project'] = moved_to_project_name_and_local_id[0]
    kwargs['moved_to_id'] = moved_to_project_name_and_local_id[1]
  elif copied_to_project_name_and_local_id:
    kwargs['copied_from_id'] = local_id
    kwargs['copied_to_project'] = copied_to_project_name_and_local_id[0]
    kwargs['copied_to_id'] = copied_to_project_name_and_local_id[1]
  else:
    kwargs['updated'] = local_id
  # if issue is being visited via a hotlist and it gets moved to another
  # project, going to issue list should mean going to hotlistissues list.
  issue_kwargs = {}
  if hotlist_id:
    url = framework_helpers.FormatAbsoluteURL(
        mr, '/u/%s/hotlists/%s' % (mr.auth.user_id, hotlist_id),
        include_project=False, **kwargs)
    issue_kwargs['hotlist_id'] = hotlist_id
  else:
    url = tracker_helpers.FormatIssueListURL(
        mr, config, **kwargs)

  if after_issue_update == user_pb2.IssueUpdateNav.STAY_SAME_ISSUE:
    # If it was a move request then will have to switch to the new project to
    # stay on the same issue.
    if moved_to_project_name_and_local_id:
      mr.project_name = moved_to_project_name_and_local_id[0]
    issue_kwargs['id'] = local_id
    url = framework_helpers.FormatAbsoluteURL(
        mr, urls.ISSUE_DETAIL_LEGACY, **issue_kwargs)
  elif after_issue_update == user_pb2.IssueUpdateNav.NEXT_IN_LIST:
    if next_id:
      issue_kwargs['id'] = next_id
      url = framework_helpers.FormatAbsoluteURL(
          mr, urls.ISSUE_DETAIL_LEGACY, project_name=next_project,
          **issue_kwargs)

  return url


# TODO(jrobbins): do we want this?
# class IssueDerivedLabelsJSON(jsonfeed.JsonFeed)


def CheckCopyIssueRequest(
    services, mr, issue, copy_selected, copy_to, errors):
  """Process the copy issue portions of the issue update form.

  Args:
    services: A Services object
    mr: commonly used info parsed from the request.
    issue: Issue protobuf for the issue being copied.
    copy_selected: True if the user selected the Copy action.
    copy_to: A project_name or url to copy this issue to or None
      if the project name wasn't sent in the form.
    errors: The errors object for this request.

    Returns:
      The project pb for the project the issue will be copy to
      or None if the copy cannot be performed. Perhaps because
      the project does not exist, in which case copy_to and
      copy_to_project will be set on the errors object. Perhaps
      the user does not have permission to copy the issue to the
      destination project, in which case the copy_to field will be
      set on the errors object.
  """
  if not copy_selected:
    return None

  if not copy_to:
    errors.copy_to = 'No destination project specified'
    errors.copy_to_project = copy_to
    return None

  copy_to_project = services.project.GetProjectByName(mr.cnxn, copy_to)
  if not copy_to_project:
    errors.copy_to = 'No such project: ' + copy_to
    errors.copy_to_project = copy_to
    return None

  # permissions enforcement
  if not servlet_helpers.CheckPermForProject(
      mr, permissions.EDIT_ISSUE, copy_to_project):
    errors.copy_to = 'You do not have permission to copy issues to project'
    errors.copy_to_project = copy_to
    return None

  elif permissions.GetRestrictions(issue):
    errors.copy_to = (
        'Issues with Restrict labels are not allowed to be copied.')
    errors.copy_to_project = ''
    return None

  return copy_to_project


def CheckMoveIssueRequest(
    services, mr, issue, move_selected, move_to, errors):
  """Process the move issue portions of the issue update form.

  Args:
    services: A Services object
    mr: commonly used info parsed from the request.
    issue: Issue protobuf for the issue being moved.
    move_selected: True if the user selected the Move action.
    move_to: A project_name or url to move this issue to or None
      if the project name wasn't sent in the form.
    errors: The errors object for this request.

    Returns:
      The project pb for the project the issue will be moved to
      or None if the move cannot be performed. Perhaps because
      the project does not exist, in which case move_to and
      move_to_project will be set on the errors object. Perhaps
      the user does not have permission to move the issue to the
      destination project, in which case the move_to field will be
      set on the errors object.
  """
  if not move_selected:
    return None

  if not move_to:
    errors.move_to = 'No destination project specified'
    errors.move_to_project = move_to
    return None

  if issue.project_name == move_to:
    errors.move_to = 'This issue is already in project ' + move_to
    errors.move_to_project = move_to
    return None

  move_to_project = services.project.GetProjectByName(mr.cnxn, move_to)
  if not move_to_project:
    errors.move_to = 'No such project: ' + move_to
    errors.move_to_project = move_to
    return None

  # permissions enforcement
  if not servlet_helpers.CheckPermForProject(
      mr, permissions.EDIT_ISSUE, move_to_project):
    errors.move_to = 'You do not have permission to move issues to project'
    errors.move_to_project = move_to
    return None

  elif permissions.GetRestrictions(issue):
    errors.move_to = (
        'Issues with Restrict labels are not allowed to be moved.')
    errors.move_to_project = ''
    return None

  return move_to_project


def _GetBinnedHotlistViews(visible_hotlist_views, involved_users):
  """Bins into (logged-in user's, issue-involved users', others') hotlists"""
  user_issue_hotlist_views = []
  involved_users_issue_hotlist_views = []
  remaining_issue_hotlist_views = []

  for view in visible_hotlist_views:
    if view.role_name in ('owner', 'editor'):
      user_issue_hotlist_views.append(view)
    elif view.owner_ids[0] in involved_users:
      involved_users_issue_hotlist_views.append(view)
    else:
      remaining_issue_hotlist_views.append(view)

  return (user_issue_hotlist_views, involved_users_issue_hotlist_views,
          remaining_issue_hotlist_views)

def _ComputeBackToListURL(mr, issue, config, hotlist, services):
  """Construct a URL to return the user to the place that they came from."""
  if hotlist:
    back_to_list_url = hotlist_helpers.GetURLOfHotlist(
        mr.cnxn, hotlist, services.user)
  else:
    back_to_list_url = tracker_helpers.FormatIssueListURL(
        mr, config, cursor='%s:%d' % (issue.project_name, issue.local_id))

  return back_to_list_url


class FlipperRedirectBase(servlet.Servlet):

  # pylint: disable=arguments-differ
  # pylint: disable=unused-argument
  def get(self, project_name=None, viewed_username=None, hotlist_id=None):
    with work_env.WorkEnv(self.mr, self.services) as we:
      hotlist_id = self.mr.GetIntParam('hotlist_id')
      current_issue = we.GetIssueByLocalID(self.mr.project_id, self.mr.local_id,
                                   use_cache=False)
      hotlist = None
      if hotlist_id:
        try:
          hotlist = self.services.features.GetHotlist(self.mr.cnxn, hotlist_id)
        except features_svc.NoSuchHotlistException:
          pass

      try:
        adj_issue = GetAdjacentIssue(we, current_issue,
            hotlist=hotlist, next_issue=self.next_handler)
        path = '/p/%s%s' % (adj_issue.project_name, urls.ISSUE_DETAIL)
        url = framework_helpers.FormatURL(
            [(name, self.mr.GetParam(name)) for
             name in framework_helpers.RECOGNIZED_PARAMS],
            path, id=adj_issue.local_id)
      except exceptions.NoSuchIssueException:
        config = we.GetProjectConfig(self.mr.project_id)
        url = _ComputeBackToListURL(self.mr, current_issue, config,
                                                 hotlist, self.services)
      self.redirect(url)


class FlipperNext(FlipperRedirectBase):
  next_handler = True


class FlipperPrev(FlipperRedirectBase):
  next_handler = False


class FlipperList(servlet.Servlet):
  # pylint: disable=arguments-differ
  # pylint: disable=unused-argument
  def get(self, project_name=None, viewed_username=None, hotlist_id=None):
    with work_env.WorkEnv(self.mr, self.services) as we:
      hotlist_id = self.mr.GetIntParam('hotlist_id')
      current_issue = we.GetIssueByLocalID(self.mr.project_id, self.mr.local_id,
                                   use_cache=False)
      hotlist = None
      if hotlist_id:
        try:
          hotlist = self.services.features.GetHotlist(self.mr.cnxn, hotlist_id)
        except features_svc.NoSuchHotlistException:
          pass

      config = we.GetProjectConfig(self.mr.project_id)

      if hotlist:
        self.mr.ComputeColSpec(hotlist)
      else:
        self.mr.ComputeColSpec(config)

      url = _ComputeBackToListURL(self.mr, current_issue, config,
                                               hotlist, self.services)
    self.redirect(url)


class FlipperIndex(jsonfeed.JsonFeed):
  """Return a JSON object of an issue's index in search.

  This is a distinct JSON endpoint because it can be expensive to compute.
  """
  CHECK_SECURITY_TOKEN = False

  def HandleRequest(self, mr):
    hotlist_id = mr.GetIntParam('hotlist_id')
    list_url = None
    with work_env.WorkEnv(mr, self.services) as we:
      if not _ShouldShowFlipper(mr, self.services):
        return {}
      issue = we.GetIssueByLocalID(mr.project_id, mr.local_id, use_cache=False)
      hotlist = None

      if hotlist_id:
        hotlist = self.services.features.GetHotlist(mr.cnxn, hotlist_id)

        if not features_bizobj.IssueIsInHotlist(hotlist, issue.issue_id):
          raise exceptions.InvalidHotlistException()

        if not permissions.CanViewHotlist(
            mr.auth.effective_ids, mr.perms, hotlist):
          raise permissions.PermissionException()

        (prev_iid, cur_index, next_iid, total_count
            ) = we.GetIssuePositionInHotlist(issue, hotlist)
      else:
        (prev_iid, cur_index, next_iid, total_count
            ) = we.FindIssuePositionInSearch(issue)

      config = we.GetProjectConfig(self.mr.project_id)

      if hotlist:
        mr.ComputeColSpec(hotlist)
      else:
        mr.ComputeColSpec(config)

      list_url = _ComputeBackToListURL(mr, issue, config, hotlist,
        self.services)

    prev_url = None
    next_url = None

    recognized_params = [(name, mr.GetParam(name)) for name in
                           framework_helpers.RECOGNIZED_PARAMS]
    if prev_iid:
      prev_issue = we.services.issue.GetIssue(mr.cnxn, prev_iid)
      path = '/p/%s%s' % (prev_issue.project_name, urls.ISSUE_DETAIL)
      prev_url = framework_helpers.FormatURL(
          recognized_params, path, id=prev_issue.local_id)

    if next_iid:
      next_issue = we.services.issue.GetIssue(mr.cnxn, next_iid)
      path = '/p/%s%s' % (next_issue.project_name, urls.ISSUE_DETAIL)
      next_url = framework_helpers.FormatURL(
          recognized_params, path, id=next_issue.local_id)

    return {
      'prev_iid': prev_iid,
      'prev_url': prev_url,
      'cur_index': cur_index,
      'next_iid': next_iid,
      'next_url': next_url,
      'list_url': list_url,
      'total_count': total_count,
    }


def _ShouldShowFlipper(mr, services):
  """Return True if we should show the flipper."""

  # Check if the user entered a specific issue ID of an existing issue.
  if tracker_constants.JUMP_RE.match(mr.query):
    return False

  # Check if the user came directly to an issue without specifying any
  # query or sort.  E.g., through crbug.com.  Generating the issue ref
  # list can be too expensive in projects that have a large number of
  # issues.  The all and open issues cans are broad queries, other
  # canned queries should be narrow enough to not need this special
  # treatment.
  if (not mr.query and not mr.sort_spec and
      mr.can in [tracker_constants.ALL_ISSUES_CAN,
                 tracker_constants.OPEN_ISSUES_CAN]):
    num_issues_in_project = services.issue.GetHighestLocalID(
        mr.cnxn, mr.project_id)
    if num_issues_in_project > settings.threshold_to_suppress_prev_next:
      return False

  return True


def GetAdjacentIssue(we, issue, hotlist=None, next_issue=False):
  """Compute next or previous issue given params of current issue.

  Args:
    we: A WorkEnv instance.
    issue: The current issue (from which to compute prev/next).
    hotlist (optional): The current hotlist.
    next_issue (bool): If True, return next, issue, else return previous issue.

  Returns:
    The adjacent issue.

  Raises:
    NoSuchIssueException when there is no adjacent issue in the list.
  """
  if hotlist:
    (prev_iid, _cur_index, next_iid, _total_count
        ) = we.GetIssuePositionInHotlist(issue, hotlist)
  else:
    (prev_iid, _cur_index, next_iid, _total_count
        ) = we.FindIssuePositionInSearch(issue)
  iid = next_iid if next_issue else prev_iid
  if iid is None:
    raise exceptions.NoSuchIssueException()
  return we.GetIssue(iid)
