// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import sinon from 'sinon';
import {MrCommentList} from './mr-comment-list.js';
import {flush} from '@polymer/polymer/lib/utils/flush.js';
import {actionType} from '../../redux/redux-mixin.js';


let element;

suite('mr-comment-list', () => {
  setup(() => {
    element = document.createElement('mr-comment-list');
    document.body.appendChild(element);
    element.comments = [
      {
        canFlag: true,
        localId: 898395,
        canDelete: true,
        projectName: 'chromium',
        commenter: {
          displayName: 'user@example.com',
          userId: '12345',
        },
        content: 'foo',
        sequenceNum: 1,
        timestamp: 1549319989,
      },
      {
        canFlag: true,
        localId: 898395,
        canDelete: true,
        projectName: 'chromium',
        commenter: {
          displayName: 'user@example.com',
          userId: '12345',
        },
        content: 'foo',
        sequenceNum: 2,
        timestamp: 1549320089,
      },
      {
        canFlag: true,
        localId: 898395,
        canDelete: true,
        projectName: 'chromium',
        commenter: {
          displayName: 'user@example.com',
          userId: '12345',
        },
        content: 'foo',
        sequenceNum: 3,
        timestamp: 1549320189,
      },
    ];

    sinon.stub(window, 'requestAnimationFrame').callsFake((func) => func());
  });

  teardown(() => {
    document.body.removeChild(element);
    element.dispatchAction({type: actionType.RESET_STATE});

    window.requestAnimationFrame.restore();
  });

  test('initializes', () => {
    assert.instanceOf(element, MrCommentList);
  });

  test('scrolls to comment', () => {
    flush();

    const commentElement = element.shadowRoot.querySelector('#c3');
    sinon.stub(element, 'showComments');
    sinon.stub(commentElement, 'scrollIntoView');

    element.focusId = 'c3';

    flush();

    assert.isTrue(element.showComments.called);
    assert.isTrue(commentElement.scrollIntoView.calledOnce);
  });

  test('scrolls to hidden comment', () => {
    flush();

    const commentElement = element.shadowRoot.querySelector('#c1');
    sinon.stub(element, 'showComments');
    sinon.stub(commentElement, 'scrollIntoView');

    element.focusId = 'c1';

    flush();

    assert.isTrue(element.showComments.called);
    assert.isTrue(commentElement.scrollIntoView.calledOnce);
  });
});