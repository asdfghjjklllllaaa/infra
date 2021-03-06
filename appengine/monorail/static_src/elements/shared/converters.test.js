// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import {displayNameToUserRef, labelStringToRef, componentStringToRef,
  issueStringToRef, issueStringToBlockingRef, issueRefToString,
  fieldNameToLabelPrefix, commentListToDescriptionList, valueToFieldValue,
} from './converters.js';

describe('displayNameToUserRef', () => {
  it('converts displayName', () => {
    assert.deepEqual(
      displayNameToUserRef('foo@bar.com'),
      {displayName: 'foo@bar.com'});
  });
  it('throws on invalid email', () => {
    assert.throws(() => displayNameToUserRef('foo'));
  });
});

describe('labelStringToRef', () => {
  it('converts label', () => {
    assert.deepEqual(labelStringToRef('foo'), {label: 'foo'});
  });
});

describe('fieldNameToLabelPrefix', () => {
  it('converts fieldName', () => {
    assert.deepEqual(fieldNameToLabelPrefix('test'), 'test-');
    assert.deepEqual(fieldNameToLabelPrefix('test-hello'), 'test-hello-');
    assert.deepEqual(fieldNameToLabelPrefix('WHATEVER'), 'whatever-');
  });
});

describe('componentStringToRef', () => {
  it('converts component', () => {
    assert.deepEqual(componentStringToRef('foo'), {path: 'foo'});
  });
});

describe('issueStringToRef', () => {
  it('converts issue default project', () => {
    assert.deepEqual(
      issueStringToRef('proj', '1234'),
      {projectName: 'proj', localId: 1234});
  });

  it('converts issue with project', () => {
    assert.deepEqual(
      issueStringToRef('proj', 'foo:1234'),
      {projectName: 'foo', localId: 1234});
  });

  it('converts external issue references', () => {
    assert.deepEqual(
      issueStringToRef('proj', 'b/123456'),
      {extIdentifier: 'b/123456'});
  });

  it('throws on invalid input', () => {
    assert.throws(() => issueStringToRef('proj', 'foo'));
  });
});

describe('issueStringToBlockingRef', () => {
  it('converts issue default project', () => {
    assert.deepEqual(
      issueStringToBlockingRef('proj', 1, '1234'),
      {projectName: 'proj', localId: 1234});
  });

  it('converts issue with project', () => {
    assert.deepEqual(
      issueStringToBlockingRef('proj', 1, 'foo:1234'),
      {projectName: 'foo', localId: 1234});
  });

  it('throws on invalid input', () => {
    assert.throws(() => issueStringToBlockingRef('proj', 1, 'foo'));
  });

  it('throws when blocking an issue on itself', () => {
    assert.throws(() => issueStringToBlockingRef('proj', 123, 'proj:123'));
    assert.throws(() => issueStringToBlockingRef('proj', 123, '123'));
  });
});

describe('issueRefToString', () => {
  it('no ref', () => {
    assert.equal(issueRefToString(), '');
  });

  it('ref with no project name', () => {
    assert.equal(
      'other:1234',
      issueRefToString({projectName: 'other', localId: 1234})
    );
  });

  it('ref with different project name', () => {
    assert.equal(
      'other:1234',
      issueRefToString({projectName: 'other', localId: 1234}, 'proj')
    );
  });

  it('ref with same project name', () => {
    assert.equal(
      '1234',
      issueRefToString({projectName: 'proj', localId: 1234}, 'proj')
    );
  });

  it('external ref', () => {
    assert.equal(
      'b/123456',
      issueRefToString({extIdentifier: 'b/123456'}, 'proj')
    );
  });
});

describe('commentListToDescriptionList', () => {
  it('empty list', () => {
    assert.deepEqual(commentListToDescriptionList(), []);
    assert.deepEqual(commentListToDescriptionList([]), []);
  });

  it('first comment is description', () => {
    assert.deepEqual(commentListToDescriptionList([
      {content: 'test'},
      {content: 'hello'},
      {content: 'world'},
    ]), [{content: 'test'}]);
  });

  it('some descriptions', () => {
    assert.deepEqual(commentListToDescriptionList([
      {content: 'test'},
      {content: 'hello', descriptionNum: 1},
      {content: 'world'},
      {content: 'this'},
      {content: 'is a'},
      {content: 'description', descriptionNum: 2},
    ]), [
      {content: 'test'},
      {content: 'hello', descriptionNum: 1},
      {content: 'description', descriptionNum: 2},
    ]);
  });
});

describe('valueToFieldValue', () => {
  it('converts field ref and value', () => {
    assert.deepEqual(valueToFieldValue(
      {fieldName: 'name', fieldId: 'id'},
      'value',
    ), {
      fieldRef: {fieldName: 'name', fieldId: 'id'},
      value: 'value',
    });
  });
});
