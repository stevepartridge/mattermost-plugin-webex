// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import ActionTypes from '../action_types';
import Constants from '../constants';
import {getConnected, notifyCanStartMeeting} from '../actions';

export function handleWebexConnected(store) {
    
    return (msg) => {
        if (!msg.data) {
            return;
        }

        if (msg.data.success) {
          store.dispatch(getConnected());
          store.dispatch(notifyCanStartMeeting());
        }
    };
}