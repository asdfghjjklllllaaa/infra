// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import '@polymer/polymer/polymer-legacy.js';
import {PolymerElement, html} from '@polymer/polymer';

/**
 * `<chops-toggle>`
 *
 * A toggle button component. This component is primarily a wrapper
 * around a native checkbox to allow easy sharing of styles.
 *
 */
export class ChopsToggle extends PolymerElement {
  static get template() {
    return html`
      <style>
        :host {
          --chops-toggle-bg: none;
          --chops-toggle-color: black;
          --chops-toggle-hover-bg: rgba(0, 0, 0, 0.3);
          --chops-toggle-focus-border: hsl(193, 82%, 63%);
          --chops-toggle-checked-bg: rgba(0, 0, 0, 0.6);
          --chops-toggle-checked-color: white;
        }
        label {
          background: var(--chops-toggle-bg);
          color: var(--chops-toggle-color);
          cursor: pointer;
          align-items: center;
          padding: 2px 4px;
          border: 1px solid #ccc;
          border-radius: var(--chops-button-radius);
        }
        input[type="checkbox"] {
          /* We need the checkbox to be hidden but still accessible. */
          opacity: 0;
          width: 0;
          height: 0;
          position: absolute;
          top: -9999;
          left: -9999;
        }
        input[type="checkbox"]:focus + label {
          /* Make sure an outline shows around this element for
           * accessibility.
           */
           box-shadow: 0 0 5px 1px var(--chops-toggle-focus-border);
        }
        input[type="checkbox"]:hover + label {
           background: var(--chops-toggle-hover-bg);
        }
        input[type="checkbox"]:checked + label {
          background: var(--chops-toggle-checked-bg);
          color: var(--chops-toggle-checked-color);
        }
      </style>
      <!-- Note: Avoiding 2-way data binding to futureproof this code
        for LitElement. -->
      <input id="checkbox" type="checkbox" checked\$="[[checked]]" on-change="_checkedChangeHandler">
      <label for="checkbox">
        <slot></slot>
      </label>
    `;
  }

  static get is() {
    return 'chops-toggle';
  }

  static get properties() {
    return {
      /**
       * Note: At the moment, this component does not manage its own
       * internal checked state. It expects its checked state to come
       * from its parent, and its parent is expected to update the
       * chops-checkbox's checked state on a change event.
       *
       * This can be generalized in the future to support multiple
       * ways of managing checked state if needed.
       **/
      checked: {
        type: Boolean,
        observer: '_checkedChange',
      },
    };
  }

  click() {
    super.click();
    this.shadowRoot.querySelector('#checkbox').click();
  }

  _checkedChangeHandler(evt) {
    this._checkedChange(evt.target.checked);
  }

  _checkedChange(checked) {
    if (checked === this.checked) return;
    const customEvent = new CustomEvent('checked-change', {
      detail: {
        checked: checked,
      },
    });
    this.dispatchEvent(customEvent);
  }
}
customElements.define(ChopsToggle.is, ChopsToggle);