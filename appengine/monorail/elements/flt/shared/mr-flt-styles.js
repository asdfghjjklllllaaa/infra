// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

const $_documentContainer = document.createElement('template');

$_documentContainer.innerHTML = `<dom-module id="mr-flt-styles">
  <template>
    <style>
      a.button {
        /* Links that look like buttons. */
        display: inline-flex;
        align-items: center;
        justify-content: center;
        text-decoration: none;
        transition: filter 0.3s ease-in-out, box-shadow 0.3s ease-in-out;
      }
      a.button:hover {
        filter: brightness(95%);
      }
      chops-button, a.button {
        box-sizing: border-box;
        border: 1px solid hsl(0, 0%, 85%);
        font-size: 12px;
        background: white;
        border-radius: 8px;
        padding: 0.25em 8px;
        margin: 0;
        margin-left: auto;
        color: hsl(227, 60%, 39%);
      }
      chops-button i.material-icons, a.button i.material-icons {
        display: block;
        margin-right: 4px;
      }
      chops-button.emphasized, a.button.emphasized {
        background: hsl(219, 70%, 51%);
        color: white;
        text-shadow: 1px 1px 3px hsla(0, 0%, 0%, 0.25);
      }
      chops-button.de-emphasized, a.button.de-emphasized {
        border-color: transparent;
      }
      /* Note: decoupling heading levels from styles is useful for
       * accessibility because styles will not always line up with semantically
       * appropriate heading levels.
       */
      .medium-heading {
        font-size: 16px;
        font-weight: normal;
        line-height: 1;
        padding: 0.5em 0;
        color: hsl(227, 60%, 39%);
        margin: 0;
        border-bottom: 1px dotted hsl(0, 0%, 83%);
      }
    </style>
  </template>
</dom-module>`;

document.head.appendChild($_documentContainer.content);