// Wrapper for AntDesign Switch that makes an api call, then displays a confirmation icon upon 
// TODO: make it more generic, maybe. This one is currently just for message visiblity updates.

import React, { useState, useEffect } from "react";
import { Switch } from "antd";
import { CheckCircleFilled, ExclamationCircleFilled } from "@ant-design/icons";
import { fetchData, UPDATE_CHAT_MESSGAE_VIZ } from "../../utils/apis";
import { MessageType } from '../../types/chat';
import { OUTCOME_TIMEOUT } from "../chat";

interface MessageToggleProps {
  isVisible: boolean;
  message: MessageType;
  setMessage: (message: MessageType) => {},
};


export default function ToggleSwitch({ isVisible, message, setMessage }: MessageToggleProps) {
  let outcomeTimeout = null;
  const [outcome, setOutcome] = useState(0);

  const { id: messageId } = message;

  const resetOutcome = () => {
    outcomeTimeout = setTimeout(() => { setOutcome(0)}, OUTCOME_TIMEOUT);
  };
  
  useEffect(() => {
    return () => {
      clearTimeout(outcomeTimeout);
    };
  });


  const updateChatMessage = async () => {
    clearTimeout(outcomeTimeout);
    setOutcome(0);
    const result = await fetchData(UPDATE_CHAT_MESSGAE_VIZ, {
      auth: true,
      method: 'POST',
      data: {
        visible: !isVisible,
        id: messageId,
      },
    });

    if (result.success && result.message === "changed") {
      setMessage({ ...message, visible: !isVisible });
      setOutcome(1);
    } else {
      setMessage({ ...message, visible: isVisible });
      setOutcome(-1);
    }
    resetOutcome();
  }


  let outcomeIcon = <CheckCircleFilled style={{ color: 'transparent' }} />;
  if (outcome) {
    outcomeIcon = outcome > 0 ? 
    <CheckCircleFilled style={{ color: 'var(--ant-success)' }} /> :
    <ExclamationCircleFilled style={{ color: 'var(--ant-warning)' }} />;
  }

  return (
    <div className="toggle-switch">
      <span className="outcome-icon">{outcomeIcon}</span>
      <Switch
        size="small"
        onChange={updateChatMessage}
        checked={isVisible}
      />
    </div>
  );
}