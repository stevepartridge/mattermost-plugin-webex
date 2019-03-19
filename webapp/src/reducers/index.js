// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {combineReducers} from 'redux';

import ActionTypes from '../action_types';

function connected(state = false, action) {
    switch (action.type) {
    case ActionTypes.AUTH_CONNECTED:
      return action.data;
    default:
      return state;
    }
}

function disconnected(state = true, action) {
    switch (action.type) {
    case ActionTypes.AUTH_CONNECTED:
      return false;
    case ActionTypes.AUTH_DISCONNECTED:
      return true;
    default:
      return state;
    }
}

function authenticated(state = {user : {}, session: {}}, action) {
    switch (action.type) {
    case ActionTypes.AUTH_CONNECTED:
      return action.data;
    case ActionTypes.AUTH_DISCONNECTED:
      return true;
    default:
      return state;
    }
}

function meetings(state = {}, action) {
  switch (action.type) {
    case ActionTypes.MEETING_CREATED:
      return state[action.data.id] = action.data;
    default:
      return state;
  }
}

const rootModalVisible = (state = false, action) => {
    switch (action.type) {
    case ActionTypes.OPEN_ROOT_MODAL:
      return true;
    case ActionTypes.MEETING_CREATED_ERROR:
      if (action.data.message === 'Webex account not connected') {
        return true;
      }      
    case ActionTypes.AUTH_CONNECTED:
    case ActionTypes.CLOSE_ROOT_MODAL:
      return false;
    default:
      return state;
    }
};

export default combineReducers({
    connected,
    disconnected,
    authenticated,
    rootModalVisible
});
