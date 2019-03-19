// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import {id as pluginId} from '../manifest';

export default {
    AUTH_CONNECTED        : pluginId + '_auth_connected',
    AUTH_DISCONNECTED     : pluginId + '_auth_disconnected',
    AUTH_REQUIRED         : pluginId + '_auth_required',
    
    MEETING_CREATED       : pluginId + '_meeting_created',
    MEETING_CREATED_ERROR : pluginId + '_meeting_created_error',
    MEETING_STARTED       : pluginId + '_meeting_started',
    MEETING_ENDED         : pluginId + '_meeting_ended',
    
    STATUS_CHANGE         : pluginId + '_status_change',
    OPEN_ROOT_MODAL       : pluginId + '_open_root_modal',
    CLOSE_ROOT_MODAL      : pluginId + '_close_root_modal',
};
