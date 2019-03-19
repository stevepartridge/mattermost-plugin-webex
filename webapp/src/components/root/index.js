// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

const {connect} = window.ReactRedux;
const {bindActionCreators} = window.Redux;
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';
import {isRootModalVisible, getAuthenticated} from '../../selectors';
import {startMeeting, closeRootModal, getOAuthConnectURL} from '../../actions';

import Root from './root.jsx';

function mapStateToProps(state, ownProps) {
  const post = ownProps.post || {};

  return {
    visible: isRootModalVisible(state),
    authenticated: getAuthenticated(state),
    state,
    theme: getTheme(state),
    ...ownProps
  };
}

function mapDispatchToProps(dispatch) {
  /* Provide actions here if needed */
  let closePopover = closeRootModal
  return {
    actions: bindActionCreators({
      getOAuthConnectURL,
      startMeeting,
      closeRootModal,
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Root);
