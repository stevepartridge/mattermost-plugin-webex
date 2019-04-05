// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {combineReducers} from 'redux';

import ActionTypes from '../action_types';

function authenticated(state = {user: {}, session: {}}, action) {
    switch (action.type) {
    case ActionTypes.AUTH_CONNECTED:
        return action.data;
    case ActionTypes.AUTH_DISCONNECTED:
        return false;
    default:
        return state;
    }
}

/*eslint no-fallthrough: 2*/
const rootModalVisible = (state = false, action) => {
    switch (action.type) {
    case ActionTypes.OPEN_ROOT_MODAL:
        return true;
    case ActionTypes.MEETING_CREATED_ERROR:
        if (action.data.message === 'Webex account not connected') {
            return true;
        }

        // falls through
    case ActionTypes.AUTH_CONNECTED:

        // falls through
    case ActionTypes.CLOSE_ROOT_MODAL:
        return false;
    default:
        return state;
    }
};

export default combineReducers({
    authenticated,
    rootModalVisible,
});
