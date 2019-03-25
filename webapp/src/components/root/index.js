// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';
import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {isRootModalVisible, getAuthenticated} from '../../selectors';
import {startMeeting, closeRootModal} from '../../actions';

import Root from './root.jsx';

function mapStateToProps(state, ownProps) {
    const post = ownProps.post || {};

    return {
        visible: isRootModalVisible(state),
        authenticated: getAuthenticated(state),
        site_url: getConfig(state).SiteURL,
        theme: getTheme(state),
        ...ownProps,
    };
}

function mapDispatchToProps(dispatch) {
    const closePopover = closeRootModal;
    return {
        actions: bindActionCreators({
            startMeeting,
            closeRootModal,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(Root);
