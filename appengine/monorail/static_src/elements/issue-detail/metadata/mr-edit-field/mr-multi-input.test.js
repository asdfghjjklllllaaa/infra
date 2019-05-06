// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import {MrMultiInput} from './mr-multi-input.js';
import {fieldTypes} from 'elements/shared/field-types.js';

let element;

// TODO(zhangtiff): Refactor mr-multi-input to not depend on
// manual DOM editing in order to avoid having to do this.
// initialValues' change triggers another asynchronous update
// after the last one, causing this test to require multiple
// update cycles.
export const initialValueUpdateComplete = async (element) => {
  await element.updateComplete; // Wait for initialValues' change to call updated().
  await element.updateComplete; // Wait for updated() to trigger reset().
  await element.updateComplete; // Wait for reset() to finish updating the DOM.
  return true;
};

suite('mr-multi-input', () => {
  setup(() => {
    element = document.createElement('mr-multi-input');
    document.body.appendChild(element);
  });

  teardown(() => {
    document.body.removeChild(element);
  });

  test('initializes', () => {
    assert.instanceOf(element, MrMultiInput);
  });

  test('input updates when initialValues change', async () => {
    element.initialValues = ['hello', 'world'];
    await initialValueUpdateComplete(element);

    assert.deepEqual(await element.getValues(), ['hello', 'world']);
  });

  test('input updates when setValues is called', async () => {
    element.initialValues = ['hello', 'world'];
    await initialValueUpdateComplete(element);

    await element.setValues(['jaunty', 'jackalope', 'jumps', 'joyously']);
    assert.deepEqual(element.getValues(),
      ['jaunty', 'jackalope', 'jumps', 'joyously']);
  });

  test('initial value does not change after user input', async () => {
    element.initialValues = ['hello'];
    await initialValueUpdateComplete(element);

    // Simulate user input.
    element.shadowRoot.querySelector('#multi1').value = 'heron';
    assert.deepEqual(element.initialValues, ['hello']);
  });

  test('resetting input to initial value works', async () => {
    element.initialValues = [];
    await initialValueUpdateComplete(element);

    // Simulate user input.
    element.shadowRoot.querySelector('#multi0').value = 'heron';
    element.reset();
    await element.updateComplete;

    assert.deepEqual(element.initialValues, []);
    assert.equal(element.shadowRoot.querySelector('#multi0').value.trim(), '');
  });

  test('get value after user input', async () => {
    element.initialValues = ['hello'];
    await initialValueUpdateComplete(element);
    // Simulate user input.
    element.shadowRoot.querySelector('#multi1').value = 'heron';
    assert.deepEqual(element.getValues(), ['hello', 'heron']);
  });

  test('input value was added', async () => {
    await element.updateComplete;
    // Simulate user input.
    element.shadowRoot.querySelector('#multi0').value = 'jackalope';
    assert.deepEqual(element.getValues(), ['jackalope']);
  });

  test('input value was removed', async () => {
    element.initialValues = ['hello'];
    await initialValueUpdateComplete(element);
    // Simulate user input.
    element.shadowRoot.querySelector('#multi0').value = '';
    assert.deepEqual(element.getValues(), []);
  });

  test('input value was changed', async () => {
    element.initialValues = ['hello'];
    await initialValueUpdateComplete(element);
    // Simulate user input.
    element.shadowRoot.querySelector('#multi0').value = 'world';
    assert.deepEqual(element.getValues(), ['world']);
  });

  test('input value has commas', async () => {
    await element.updateComplete;
    element.acType = 'member';

    // Simulate user input.
    const input = element.shadowRoot.querySelector('#multi0');
    input.value = 'jaunty;jackalope,, jumps joyously!';
    await element._onBlur();

    await element.updateComplete;

    // Input is split on several input fields.
    assert.deepEqual(
      Array.from(element.shadowRoot.querySelectorAll('input')).map(
        (input) => input.value),
      ['jaunty', 'jackalope', 'jumps', 'joyously!', '']);

    // values are updated
    assert.deepEqual(
      element.getValues(), ['jaunty', 'jackalope', 'jumps', 'joyously!']);
  });

  test('input value has commas but is not delimitable', async () => {
    await element.updateComplete;
    element.type = fieldTypes.STR_TYPE;

    // Simulate user input.
    const input = element.shadowRoot.querySelector('#multi0');
    input.value = 'jaunty;jackalope,, jumps joyously!';

    await element._onBlur();
    await element.updateComplete;

    // Input is not split into several input fields.
    assert.deepEqual(
      Array.from(element.shadowRoot.querySelectorAll('input')).map(
        (input) => input.value),
      ['jaunty;jackalope,, jumps joyously!', '']);
    assert.deepEqual(
      element.getValues(), ['jaunty;jackalope,, jumps joyously!']);

    element.type = fieldTypes.DATE_TYPE;
    await element._onBlur();
    await element.updateComplete;

    // Input is not split into several input fields.
    assert.deepEqual(
      Array.from(element.shadowRoot.querySelectorAll('input')).map(
        (input) => input.value),
      ['jaunty;jackalope,, jumps joyously!', '']);
    assert.deepEqual(
      element.getValues(), ['jaunty;jackalope,, jumps joyously!']);

    element.type = fieldTypes.URL_TYPE;
    await element._onBlur();
    await element.updateComplete;

    // Input is not split into several input fields.
    assert.deepEqual(
      Array.from(element.shadowRoot.querySelectorAll('input')).map(
        (input) => input.value),
      ['jaunty;jackalope,, jumps joyously!', '']);
    assert.deepEqual(
      element.getValues(), ['jaunty;jackalope,, jumps joyously!']);
  });
});