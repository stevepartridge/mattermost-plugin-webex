// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';

import {id as pluginId} from './manifest';

import Icon from './components/icon.jsx';
import PostTypeWebex from './components/post_type_webex';
import Root from './components/root';
import {getConnected, startMeeting} from './actions';
import {handleWebexConnected} from './websocket';

import Reducer from './reducers';

require.context('./external/', true, /\.(js|css)$/);

class Plugin {
    // eslint-disable-next-line no-unused-vars
    async initialize(registry, store) {
        registry.registerReducer(Reducer);

        await getConnected()(store.dispatch, store.getState);

        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            (channel) => {
                startMeeting(channel.id)(store.dispatch, store.getState);
            },
            'Start Webex Meeting'
        );
        registry.registerPostTypeComponent('custom_webex', PostTypeWebex);

        registry.registerWebSocketEventHandler('custom_webex_oauth_success', handleWebexConnected(store));

        registry.registerRootComponent(Root);
    }
}

window.registerPlugin(pluginId, new Plugin());
