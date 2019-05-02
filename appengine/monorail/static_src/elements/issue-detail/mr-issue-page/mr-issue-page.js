// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import '@polymer/polymer/polymer-legacy.js';
import {PolymerElement, html} from '@polymer/polymer';

import '../../chops/chops-button/chops-button.js';
import './mr-issue-header.js';
import './mr-restriction-indicator';
import '../mr-issue-details/mr-issue-details.js';
import '../mr-metadata/mr-issue-metadata.js';
import '../mr-launch-overview/mr-launch-overview.js';
import {store, connectStore} from 'elements/reducers/base.js';
import * as issue from 'elements/reducers/issue.js';
import * as project from 'elements/reducers/project.js';
import * as user from 'elements/reducers/user.js';
import 'elements/shared/mr-shared-styles.js';

import '../dialogs/mr-edit-description/mr-edit-description.js';
import '../dialogs/mr-move-copy-issue/mr-move-copy-issue.js';
import '../dialogs/mr-convert-issue/mr-convert-issue.js';
import '../dialogs/mr-related-issues-table/mr-related-issues-table.js';
import '../dialogs/mr-update-issue-hotlists/mr-update-issue-hotlists.js';
import '../dialogs/mr-cues/mr-cues.js';

const APPROVAL_COMMENT_COUNT = 5;
const DETAIL_COMMENT_COUNT = 100;

/**
 * `<mr-issue-page>`
 *
 * The main entry point for a Monorail issue detail page.
 *
 */
export class MrIssuePage extends connectStore(PolymerElement) {
  static get template() {
    return html`
      <style include="mr-shared-styles">
        :host {
          --mr-issue-page-horizontal-padding: 12px;
          --mr-toggled-font-family: inherit;
          --monorail-metadata-toggled-bg: var(--monorail-metadata-open-bg);
        }
        :host([issue-closed]) {
          --monorail-metadata-toggled-bg: var(--monorail-metadata-closed-bg);
        }
        :host([code-font]) {
          --mr-toggled-font-family: Monospace;
        }
        .container-issue {
          width: 100%;
          flex-direction: column;
          align-items: stretch;
          justify-content: flex-start;
        }
        .container-issue-content {
          padding: 0;
          overflow: hidden;
          flex-grow: 1;
          display: flex;
          align-items: stretch;
          justify-content: space-between;
          flex-direction: row;
          flex-wrap: nowrap;
          box-sizing: border-box;
          padding-top: 0.5em;
        }
        .container-outside {
          box-sizing: border-box;
          width: 100%;
          max-width: 100%;
          margin: auto;
          padding: 0;
          display: flex;
          align-items: stretch;
          justify-content: space-between;
          flex-direction: row;
          flex-wrap: no-wrap;
        }
        .container-no-issue {
          padding: 0.5em 16px;
          font-size: var(--chops-large-font-size);
        }
        .metadata-container {
          font-size: var(--chops-main-font-size);
          background: var(--monorail-metadata-toggled-bg);
          border-right: var(--chops-normal-border);
          border-bottom: var(--chops-normal-border);
          width: 16%;
          min-width: 256px;
          flex-grow: 0;
          flex-shrink: 0;
          box-sizing: border-box;
          z-index: 100;
        }
        .issue-header-container {
          z-index: 10;
          position: sticky;
          top: var(--monorail-header-height);
          margin-bottom: 0.25em;
          width: 100%;
        }
        mr-issue-details {
          min-width: 50%;
          max-width: 1000px;
          flex-grow: 1;
          box-sizing: border-box;
          min-height: 100%;
          padding-left: var(--mr-issue-page-horizontal-padding);
          padding-right: var(--mr-issue-page-horizontal-padding);
        }
        mr-issue-metadata {
          position: sticky;
          top: var(--monorail-header-height);
        }
        mr-launch-overview {
          border-left: var(--chops-normal-border);
          padding-left: var(--mr-issue-page-horizontal-padding);
          padding-right: var(--mr-issue-page-horizontal-padding);
          flex-grow: 0;
          flex-shrink: 0;
          width: 50%;
          box-sizing: border-box;
          min-height: 100%;
        }
        @media (max-width: 1024px) {
          .container-issue-content {
            flex-direction: column;
            padding: 0 var(--mr-issue-page-horizontal-padding);
          }
          mr-issue-details, mr-launch-overview {
            width: 100%;
            padding: 0;
            border: 0;
          }
        }
        @media (max-width: 840px) {
          .container-outside {
            flex-direction: column;
          }
          .metadata-container {
            width: 100%;
            height: auto;
            border: 0;
            border-bottom: var(--chops-normal-border);
          }
          mr-issue-metadata {
            min-width: auto;
            max-width: auto;
            width: 100%;
            padding: 0;
            min-height: 0;
            border: 0;
          }
          mr-issue-metadata, .issue-header-container {
            position: static;
          }
        }
      </style>
      <mr-cues id="cues" user-display-name="[[userDisplayName]]"></mr-cues>
      <template is="dom-if" if="[[_showLoading(fetchingIssue, issue)]]">
        <div class="container-no-issue" id="loading">
          Loading...
        </div>
      </template>
      <template is="dom-if" if="[[!_showLoading(fetchingIssue, issue)]]">
        <template is="dom-if" if="[[fetchIssueError]]">
          <div class="container-no-issue" id="fetch-error">
            [[fetchIssueError.description]]
          </div>
        </template>
        <template is="dom-if" if="[[_issueIsDeleted(issue)]]">
          <div class="container-no-issue" id="deleted">
            <p>Issue [[issueRef.localId]] has been deleted.</p>
            <template is="dom-if" if="[[_showUndelete(issuePermissions)]]">
              <chops-button on-click="_undeleteIssue" class="undelete emphasized">
                Undelete Issue
              </chops-button>
            </template>
          </div>
        </template>
        <template is="dom-if" if="[[_issueIsMoved(issue.movedToRef)]]">
          <div class="container-no-issue" id="moved">
            <h2>Issue has moved.</h2>
            <p>
              This issue was moved to [[issue.movedToRef.projectName]].
              <a class="new-location" href$="/p/[[issue.movedToRef.projectName]]/issues/detail?id=[[issue.movedToRef.localId]]">
                Go to issue</a>.
            </p>
          </div>
        </template>
        <template is="dom-if" if="[[!_issueIsEmpty(issue)]]">
          <div
            class="container-outside"
            on-open-dialog="_onOpenDialog"
            id="issue"
          >
            <aside class="metadata-container">
              <mr-issue-metadata></mr-issue-metadata>
            </aside>
            <div class="container-issue">
              <div class="issue-header-container">
                <mr-issue-header
                  user-display-name="[[userDisplayName]]"
                ></mr-issue-header>
                <mr-restriction-indicator></mr-restriction-indicator>
              </div>
              <div class="container-issue-content">
                <mr-issue-details
                  class="main-item"
                  comments-shown-count="[[_commentsShownCount(issue)]]"
                ></mr-issue-details>
                <mr-launch-overview class="main-item"></mr-launch-overview>
              </div>
            </div>
          </div>
          <mr-edit-description id="edit-description"></mr-edit-description>
          <mr-move-copy-issue id="move-copy-issue"></mr-move-copy-issue>
          <mr-convert-issue id="convert-issue"></mr-convert-issue>
          <mr-related-issues-table id="reorder-related-issues"></mr-related-issues-table>
          <mr-update-issue-hotlists id="update-issue-hotlists"></mr-update-issue-hotlists>
        </template>
      </template>
    `;
  }

  static get is() {
    return 'mr-issue-page';
  }

  static get properties() {
    return {
      issueEntryUrl: String,
      queryParams: {
        type: Object,
        value: () => {},
      },
      userDisplayName: String,
      // Redux state.
      fetchIssueError: String,
      fetchingIssue: Boolean,
      fetchingProjectConfig: Boolean,
      issue: Object,
      issueClosed: {
        type: Boolean,
        reflectToAttribute: true,
      },
      issuePermissions: Object,
      issueRef: Object,
      prefs: Object,
      codeFont: {
        type: Boolean,
        computed: '_computeCodeFont(prefs)',
        reflectToAttribute: true,
      },
    };
  }

  stateChanged(state) {
    this.setProperties({
      issue: issue.issue(state),
      issueRef: issue.issueRef(state),
      fetchIssueError: issue.requests(state).fetch.error,
      fetchingIssue: issue.requests(state).fetch.requesting,
      fetchingProjectConfig: project.fetchingConfig(state),
      issueClosed: !issue.isOpen(state),
      issuePermissions: issue.permissions(state),
      prefs: user.user(state).prefs,
    });
  }

  static get observers() {
    return [
      '_issueChanged(issueRef, issue)',
      '_fetchIssue(issueRef)',
      '_projectNameChanged(issueRef.projectName)',
      '_userDisplayNameChanged(userDisplayName)',
    ];
  }

  _issueChanged(issueRef, issue) {
    if (!issueRef) return;
    let title =
      issueRef.localId ? `${issueRef.localId} - ` : 'Loading issue... - ';
    if (issue && issue.isDeleted) {
      title += 'Issue has been deleted - ';
    } else if (issue && issue.summary) {
      title += `${issue.summary} - `;
    }
    if (issueRef.projectName) {
      title += `${issueRef.projectName} - `;
    }
    title += 'Monorail';
    document.title = title;
  }

  _fetchIssue(issueRef) {
    if (issueRef.localId && issueRef.projectName && !this.fetchingIssue) {
      // Reload the issue data when the id changes.
      store.dispatch(issue.fetchIssuePageData({issueRef}));
    }
  }

  _projectNameChanged(projectName) {
    if (projectName && !this.fetchingProjectConfig) {
      // Reload project config and templates when the project name changes.
      store.dispatch(project.fetch(projectName));
    }
  }

  _userDisplayNameChanged(userDisplayName) {
    store.dispatch(user.fetch(userDisplayName));
  }

  _onOpenDialog(e) {
    this.shadowRoot.querySelector('#' + e.detail.dialogId).open(e);
  }

  _commentsShownCount(issue) {
    return issue.approvalValues ? APPROVAL_COMMENT_COUNT : DETAIL_COMMENT_COUNT;
  }

  _undeleteIssue() {
    window.prpcClient.call('monorail.Issues', 'DeleteIssue', {
      issueRef: this.issueRef,
      delete: false,
    }).then(() => {
      this._fetchIssue(this.issueRef);
    });
  }

  _showUndelete(issuePermissions) {
    return (issuePermissions || []).includes('deleteissue');
  }

  _computeCodeFont(prefs) {
    if (!prefs) return false;
    return prefs.get('code_font') === 'true';
  }

  _issueIsEmpty(issue) {
    return !issue || !issue.localId;
  }

  _issueIsDeleted(issue) {
    return issue && issue.isDeleted;
  }

  _issueIsMoved(ref) {
    return ref && ref.projectName && ref.localId;
  }

  _showLoading(fetchingIssue, issue) {
    return fetchingIssue && this._issueIsEmpty(issue);
  }
}

customElements.define(MrIssuePage.is, MrIssuePage);