// Session recovery dialog component
import React, { useEffect, useState } from 'react';
import { useFormContext } from '../../contexts/FormContext';
import Button from '../common/Button';
import { SESSION_CONFIG } from '../../utils/constants';

interface SessionRecoveryDialogProps {
  onRecover: () => void;
  onDiscard: () => void;
  onClose: () => void;
}

const SessionRecoveryDialog: React.FC<SessionRecoveryDialogProps> = ({
  onRecover,
  onDiscard,
  onClose
}) => {
  const [sessionAge, setSessionAge] = useState<string>('');

  useEffect(() => {
    const lastSaved = localStorage.getItem(SESSION_CONFIG.STORAGE_KEYS.LAST_SAVED);
    if (lastSaved) {
      const lastSavedTime = new Date(lastSaved);
      const now = new Date();
      const diffMinutes = Math.floor((now.getTime() - lastSavedTime.getTime()) / (1000 * 60));
      
      if (diffMinutes < 60) {
        setSessionAge(`${diffMinutes}分前`);
      } else if (diffMinutes < 1440) {
        const hours = Math.floor(diffMinutes / 60);
        setSessionAge(`${hours}時間前`);
      } else {
        const days = Math.floor(diffMinutes / 1440);
        setSessionAge(`${days}日前`);
      }
    }
  }, []);

  const handleRecover = () => {
    onRecover();
    onClose();
  };

  const handleDiscard = () => {
    onDiscard();
    onClose();
  };

  return (
    <div className="session-recovery-dialog" role="dialog" aria-labelledby="recovery-title" aria-describedby="recovery-message">
      <div className="session-recovery-dialog__overlay" onClick={onClose} />
      <div className="session-recovery-dialog__modal">
        <div className="session-recovery-dialog__header">
          <h3 id="recovery-title" className="session-recovery-dialog__title">
            保存されたデータが見つかりました
          </h3>
        </div>
        
        <div className="session-recovery-dialog__content">
          <p id="recovery-message" className="session-recovery-dialog__message">
            前回の入力内容が残っています。続きから入力しますか？
          </p>
          
          {sessionAge && (
            <div className="session-recovery-dialog__info">
              <span className="session-recovery-dialog__timestamp">
                最後の保存: {sessionAge}
              </span>
            </div>
          )}
          
          <div className="session-recovery-dialog__warning">
            <h4>ご注意</h4>
            <ul>
              <li>「続きから入力」を選択すると、前回の入力内容が復元されます</li>
              <li>「最初から入力」を選択すると、保存されたデータは削除されます</li>
              <li>データの復元後も、内容の確認・修正が可能です</li>
            </ul>
          </div>
        </div>
        
        <div className="session-recovery-dialog__actions">
          <Button
            variant="primary"
            onClick={handleRecover}
            data-testid="recover-session"
          >
            続きから入力
          </Button>
          
          <Button
            variant="secondary"
            onClick={handleDiscard}
            data-testid="discard-session"
          >
            最初から入力
          </Button>
          
          <Button
            variant="outline"
            onClick={onClose}
            data-testid="cancel-recovery"
          >
            後で決める
          </Button>
        </div>
      </div>
    </div>
  );
};

export default SessionRecoveryDialog;