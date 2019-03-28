// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {id} from '../manifest';

import {unauthorized} from '../actions';

export default class Client {
    constructor() {
        this.url = '/plugins/' + id;
    }

    getConnected = async () => {
        const req = this.request('GET', '/connected');

        try {
            const response = await fetch(req);
            return await this.response(response);
        } catch (e) {
            return e.message;
        }
    }

    startMeeting = async (channelId, personal = true, topic = '', meetingId = 0) => {
        const payload = {
            channel_id: channelId,
            personal,
            topic,
            meeting_id: meetingId,
        };

        const req = this.request('POST', '/api/v1/meetings', payload);

        try {
            const response = await fetch(req);
            return await this.response(response);
        } catch (e) {
            return e.message;
        }
    }

    // getWebexUser = async (userID) => {
    //     return this.doPost(`${this.url}/api/v1/user`, {user_id: userID});
    // }

    request = (method, path, payload) => {
        var req = {
            method,
            headers: {
                'X-Requested-With': 'XMLHttpRequest',
                'X-Timezone-Offset': new Date().getTimezoneOffset(),
            },
        };

        if (payload) {
            req.body = JSON.stringify(payload);
        }

        const r = new Request(`${this.url}${path}`,
            req
        );

        return r;
    }

    response = (response) => {
        switch (response.status) {
        case 403:
            // permissionDenied();
            break;
        case 401:
            unauthorized();
            break;
        default:
            // success();
            break;
        }

        return response.json().then((data) => (
            data
        )
        );
    }
}
