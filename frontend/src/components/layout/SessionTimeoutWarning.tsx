// Enhanced session timeout warning component
import React from 'react';
import { useSessionManager } from '../../hooks/useSessionManager';
import Button from '../common/Button';
import { BUTTON_LABELS } from '../../utils/constants';

const SessionTimeoutWarning: React.FC = () => {
  const { 
    timeoutWarning, 
    extendSession, 
    clearSession,
    saveSession 
  } = useSessionManager();

  const handleExtendSession = async () => {
    await extendSession();
  };

  const handleSaveAndExit = async () => {
    await saveSession();
    window.location.href = '/';
  };

  const handleDiscardAndExit = async () => {
    await clearSession();
    window.location.href = '/';
  };

  if (!timeoutWarning.show) {
    return null;
  }

  return (
    <div className="session-timeout-warning" role="dialog" aria-labelledby="timeout-title" aria-describedby="timeout-message">
      <div className="session-timeout-warning__overlay" />
      <div className="session-timeout-warning__modal">
        <div className="session-timeout-warning__header">
          <h3 id="timeout-title" className="session-timeout-warning__title">
            セッション期限が近づいています
          </h3>
        </div>
        
        <div className="session-timeout-warning__content">
          <p id="timeout-message" className="session-timeout-warning__message">
            入力内容が消去される前に、セッションを延長するか保存してください。
          </p>
          <p className="session-timeout-warning__time">
            残り時間: <strong>{timeoutWarning.remainingMinutes}分</strong>
          </p>
          
          <div className="session-timeout-warning__options">
            <h4>選択してください：</h4>
            <ul>
              <li><strong>セッション延長：</strong>4時間延長して入力を続ける</li>
              <li><strong>保存して終了：</strong>現在の内容を保存して後で再開</li>
              <li><strong>破棄して終了：</strong>入力内容を破棄して最初から</li>
            </ul>
          </div>
        </div>
        
        <div className="session-timeout-warning__actions">
          <Button
            variant="primary"
            onClick={handleExtendSession}
            data-testid="extend-session-button"
          >
            {BUTTON_LABELS.EXTEND_SESSION}
          </Button>
          
          <Button
            variant="secondary"
            onClick={handleSaveAndExit}
            data-testid="save-and-exit-button"
          >
            保存して終了
          </Button>
          
          <Button
            variant="outline"
            onClick={handleDiscardAndExit}
            data-testid="discard-and-exit-button"
          >
            破棄して終了
          </Button>
        </div>
      </div>
    </div>
  );
};

export default SessionTimeoutWarning;