// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import '@polymer/polymer/polymer-legacy.js';
import {PolymerElement, html} from '@polymer/polymer';

import '../../chops/chops-dialog/chops-dialog.js';
import {selectors} from '../../redux/selectors.js';
import {actionCreator} from '../../redux/redux-mixin.js';
import '../mr-approval-card/mr-approval-card.js';
import '../mr-edit-metadata/mr-edit-metadata.js';
import '../mr-metadata/mr-field-values.js';
import {MetadataMixin} from '../shared/metadata-mixin.js';
import '../shared/mr-flt-styles.js';

const TARGET_PHASE_MILESTONE_MAP = {
  'Beta': 'feature_freeze',
  'Stable-Exp': 'final_beta_cut',
  'Stable': 'stable_cut',
  'Stable-Full': 'stable_cut',
};

const APPROVED_PHASE_MILESTONE_MAP = {
  'Beta': 'earliest_beta',
  'Stable-Exp': 'final_beta',
  'Stable': 'stable_date',
  'Stable-Full': 'stable_date',
};

// See monorail:4692 and the use of PHASES_WITH_MILESTONES
// in tracker/issueentry.py
const PHASES_WITH_MILESTONES = ['Beta', 'Stable', 'Stable-Exp', 'Stable-Full'];

/**
 * `<mr-phase>`
 *
 * This is the component for a single phase.
 *
 */
export class MrPhase extends MetadataMixin(PolymerElement) {
  static get template() {
    return html`
      <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
      <style include="mr-flt-styles">
        :host {
          display: block;
        }
        chops-dialog {
          font-size: 85%;
          --chops-dialog-theme: {
            width: 500px;
            max-width: 100%;
          };
        }
        h2 {
          margin: 0;
          font-size: 105%;
          font-weight: normal;
          padding: 0.5em 8px;
          box-sizing: border-box;
          display: flex;
          align-items: center;
          flex-direction: row;
          justify-content: space-between;
        }
        h2 em {
          margin-left: 16px;
          font-size: 80%;
        }
        .phase-edit {
          padding: 0.1em 8px;
        }
      </style>
      <h2>
        <span>
          Approvals<span hidden\$="[[_isEmpty(phaseName)]]">:
            [[phaseName]]
          </span>
          <span hidden\$="[[!_setPhaseFields.length]]">
            -
            <template is="dom-repeat" items="[[_setPhaseFields]]" as="field">
              [[field.fieldRef.fieldName]]
              <mr-field-values
                name="[[field.fieldRef.fieldName]]"
                type="[[field.fieldRef.type]]"
                values="[[_valuesForField(fieldValueMap, field.fieldRef.fieldName, phaseName)]]"
                project-name="[[projectName]]"
              ></mr-field-values><span hidden\$="[[_isLastItem(_setPhaseFields.length, index)]]">;
            </span></template>
            <em hidden\$="[[!_nextDate]]">
              [[_dateDescriptor]]
              <chops-timestamp timestamp="[[_nextDate]]"></chops-timestamp>
            </em>
          </span>
        </span>
        <chops-button hidden\$="[[_isEmpty(phaseName)]]" on-click="edit" class="de-emphasized phase-edit">
          <i class="material-icons">create</i>
          Edit
        </chops-button>
      </h2>
      <template is="dom-repeat" items="[[approvals]]">
        <mr-approval-card approvers="[[item.approverRefs]]" setter="[[item.setterRef]]" field-name="[[item.fieldRef.fieldName]]" phase-name="[[phaseName]]" status-enum="[[item.status]]" survey="[[item.survey]]" survey-template="[[item.surveyTemplate]]" urls="[[item.urls]]" labels="[[item.labels]]" users="[[item.users]]"></mr-approval-card>
      </template>
      <template is="dom-if" if="[[!approvals.length]]">
        No tasks for this phase.
      </template>
      <chops-dialog id="editPhase" aria-labelledby="phaseDialogTitle">
        <h3 id="phaseDialogTitle" class="medium-heading">
          Editing phase: [[phaseName]]
        </h3>
        <mr-edit-metadata id="metadataForm" field-defs="[[fieldDefs]]" field-values="[[fieldValues]]" phase-name="[[phaseName]]" disabled="[[updatingIssue]]" error="[[updateIssueError.description]]" on-save="save" on-discard="cancel" is-approval="" disable-attachments=""></mr-edit-metadata>
      </chops-dialog>
    `;
  }

  static get is() {
    return 'mr-phase';
  }

  static get properties() {
    return {
      issue: {
        type: Object,
        observer: 'reset',
      },
      projectName: String,
      issueId: Number,
      phaseName: {
        type: String,
        value: '',
      },
      updatingIssue: {
        type: Boolean,
        observer: '_updatingIssueChanged',
      },
      updateIssueError: Object,
      // Possible values: Target, Approved, Launched.
      _status: {
        type: String,
        computed: `_computeStatus(_targetMilestone, _approvedMilestone,
          _launchedMilestone)`,
      },
      _approvedMilestone: {
        type: Number,
        computed: '_computeApprovedMilestone(fieldValueMap, phaseName)',
      },
      _launchedMilestone: {
        type: Number,
        computed: '_computeLaunchedMilestone(fieldValueMap, phaseName)',
      },
      _targetMilestone: {
        type: Number,
        computed: '_computeTargetMilestone(fieldValueMap, phaseName)',
      },
      _fetchedMilestone: {
        type: Number,
        computed: `_computeFetchedMilestone(_targetMilestone,
          _approvedMilestone, _launchedMilestone)`,
        observer: '_fetchMilestoneData',
      },
      approvals: Array,
      fieldDefs: Array,
      fieldValueMap: {
        type: Object,
        value: () => {},
      },
      _nextDate: {
        type: Number, // Unix time.
        computed: `_computeNextDate(
          phaseName, _status, _milestoneData.mstones)`,
        value: 0,
      },
      _dateDescriptor: {
        type: String,
        computed: '_computeDateDescriptor(_status)',
      },
      _setPhaseFields: {
        type: Array,
        computed: '_computeSetPhaseFields(fieldDefs, fieldValueMap, phaseName)',
      },
      _milestoneData: Object,
    };
  }

  static mapStateToProps(state, element) {
    return {
      issue: state.issue,
      issueId: state.issueId,
      projectName: state.projectName,
      updatingIssue: state.updatingIssue,
      updateIssueError: state.updateIssueError,
      fieldDefs: selectors.fieldDefsForPhases(state),
      fieldValueMap: selectors.issueFieldValueMap(state),
    };
  }

  edit() {
    this.reset();
    this.$.editPhase.open();
  }

  cancel() {
    this.$.editPhase.close();
  }

  reset() {
    this.$.metadataForm.reset();
  }

  save() {
    const metadata = this.$.metadataForm;
    const data = metadata.getDelta();
    let message = {
      issueRef: {
        projectName: this.projectName,
        localId: this.issueId,
      },
    };

    if (metadata.sendEmail) {
      message.sendEmail = true;
    }

    let delta = {};

    const fieldValuesAdded = data.fieldValuesAdded || [];
    const fieldValuesRemoved = data.fieldValuesRemoved || [];

    if (fieldValuesAdded.length) {
      delta['fieldValsAdd'] = data.fieldValuesAdded.map(
        (fv) => Object.assign({phaseRef: {phaseName: this.phaseName}}, fv));
    }

    if (fieldValuesRemoved.length) {
      delta['fieldValsRemove'] = data.fieldValuesRemoved.map(
        (fv) => Object.assign({phaseRef: {phaseName: this.phaseName}}, fv));
    }

    if (data.comment) {
      message['commentContent'] = data.comment;
    }

    if (Object.keys(delta).length > 0) {
      message.delta = delta;
    }

    if (message.commentContent || message.delta) {
      actionCreator.updateIssue(this.dispatchAction.bind(this), message);
    }
  }

  _updatingIssueChanged(isUpdateInFlight) {
    if (!isUpdateInFlight && !this.updateIssueError) {
      // Close phase edit modal only after request finishes without errors.
      this.cancel();
    }
  }

  _computeNextDate(phaseName, status, data) {
    // Data pulled from https://chromepmo.appspot.com/schedule/mstone/json?mstone=xx
    if (!phaseName || !status || !data || !data.length) return 0;
    data = data[0];

    let key = TARGET_PHASE_MILESTONE_MAP[phaseName];
    if (['Approved', 'Launched'].includes(status)) {
      key = APPROVED_PHASE_MILESTONE_MAP[phaseName];
    }
    if (!key || !(key in data)) return 0;
    return Math.floor((new Date(data[key])).getTime() / 1000);
  }

  _computeDateDescriptor(status) {
    if (status === 'Approved') {
      return 'Launching on ';
    } else if (status === 'Launched') {
      return 'Launched on ';
    }
    return 'Due by ';
  }

  _computeSetPhaseFields(fieldDefs, fieldValueMap, phaseName) {
    // monorail:4692, remove later
    if (!PHASES_WITH_MILESTONES.includes(phaseName)) return [];
    if (!fieldDefs || !fieldValueMap) return [];
    return fieldDefs.filter((fd) => fieldValueMap.has(
      this._makeFieldValueMapKey(fd.fieldRef.fieldName, phaseName)
    ));
  }

  _computeMilestoneFieldValue(fieldValueMap, phaseName, fieldName) {
    const values = this._valuesForField(fieldValueMap, fieldName, phaseName);
    return values.length ? values[0] : undefined;
  }

  _computeApprovedMilestone(fieldValueMap, phaseName) {
    return this._computeMilestoneFieldValue(fieldValueMap, phaseName,
      'M-Approved');
  }

  _computeLaunchedMilestone(fieldValueMap, phaseName) {
    return this._computeMilestoneFieldValue(fieldValueMap, phaseName,
      'M-Launched');
  }

  _computeTargetMilestone(fieldValueMap, phaseName) {
    return this._computeMilestoneFieldValue(fieldValueMap, phaseName,
      'M-Target');
  }

  _computeStatus(target, approved, launched) {
    if (approved >= target) {
      if (launched >= approved) {
        return 'Launched';
      }
      return 'Approved';
    }
    return 'Target';
  }

  _computeFetchedMilestone(target, approved, launched) {
    return Math.max(target || 0, approved || 0, launched || 0);
  }

  _fetchMilestoneData(milestone) {
    if (!milestone) return;
    // HACK. Eventually we want to create a less bespoke way of getting
    // milestone metadata into Monorail.
    window.fetch(
      `https://chromepmo.appspot.com/schedule/mstone/json?mstone=${milestone}`
    ).then((resp) => resp.json()).then((resp) => {
      this._milestoneData = resp;
    });
  }

  _isEmpty(str) {
    return !str || !str.length;
  }

  _isLastItem(l, i) {
    return i >= l - 1;
  }
}
customElements.define(MrPhase.is, MrPhase);