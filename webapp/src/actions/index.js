// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {PostTypes} from 'mattermost-redux/action_types';

import Client from '../client';
import ActionTypes from '../action_types';

export const openRootModal = () => (dispatch) => {
    dispatch({
        type: ActionTypes.OPEN_ROOT_MODAL,
    });
};

export const closeRootModal = () => (dispatch) => {
    dispatch({
        type: ActionTypes.CLOSE_ROOT_MODAL,
    });
};

export function unauthorized() {
    return (dispatch) => {
        dispatch({
            type: ActionTypes.AUTH_DISCONNECTED,
        });
    };
}

export function getConnected() {
    return (dispatch) => {
        return Client.getConnected().then((response) => {
            if (!response) {
                dispatch({
                    type: ActionTypes.AUTH_DISCONNECTED,
                });
                return;
            }
            dispatch({
                type: ActionTypes.AUTH_CONNECTED,
                data: response,
            });

            return {response};
        }).catch((error) => {
            throw (error);
        });
    };
}

export function startMeeting(channelId) {
    return async (dispatch, getState) => {
        console.log('start meeting');

        let data;
        try {
            data = await Client.startMeeting(channelId, true);
        } catch (error) {
            return {error};
        }

        console.log('data', data);

        if (data.error) {
            dispatch({
                type: ActionTypes.MEETING_CREATED_ERROR,
                data,
            });

            const post = {
                id: 'webexPlugin' + Date.now(),
                create_at: Date.now(),
                update_at: 0,
                edit_at: 0,
                delete_at: 0,
                is_pinned: false,
                user_id: getState().entities.users.currentUserId,
                channel_id: channelId,
                root_id: '',
                parent_id: '',
                original_id: '',
                message: data.message,
                type: 'system_ephemeral',
                props: {},
                hashtags: '',
                pending_post_id: '',
            };

            dispatch({
                type: PostTypes.RECEIVED_POSTS,
                data: {
                    order: [],
                    posts: {
                        [post.id]: post,
                    },
                },
                channelId,
            });

            return data;
        }

        dispatch({
            type: ActionTypes.MEETING_CREATED,
            data,
        });

        return {data};
    };
}

export function notifyCanStartMeeting() {
    return (dispatch, getState) => {
        const channelId = getState().entities.channels.currentChannelId;

        const post = {
            id: 'webexPlugin' + Date.now(),
            create_at: Date.now(),
            update_at: 0,
            edit_at: 0,
            delete_at: 0,
            is_pinned: false,
            user_id: getState().entities.users.currentUserId,
            channel_id: channelId,
            root_id: '',
            parent_id: '',
            original_id: '',
            message: 'Webex account connected.  You can now start a meeting.',
            type: 'system_ephemeral',
            props: {},
            hashtags: '',
            pending_post_id: '',
        };

        dispatch({
            type: PostTypes.RECEIVED_POSTS,
            data: {
                order: [],
                posts: {
                    [post.id]: post,
                },
            },
            channelId,
        });
    };
}
