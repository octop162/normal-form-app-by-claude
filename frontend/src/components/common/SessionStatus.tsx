// Session status indicator component
import React from 'react';
import { useSessionManager } from '../../hooks/useSessionManager';
import Button from './Button';
import LoadingSpinner from './LoadingSpinner';

interface SessionStatusProps {
  showDetails?: boolean;
  className?: string;
  'data-testid'?: string;
}

const SessionStatus: React.FC<SessionStatusProps> = ({
  showDetails = false,
  className = '',
  'data-testid': testId
}) => {
  const {
    sessionExists,
    isSessionLoading,
    sessionError,
    saveSession,
    clearSession,
    extendSession,
    timeoutWarning,
    isAutoSaveEnabled,
    enableAutoSave,
    disableAutoSave
  } = useSessionManager();

  const handleManualSave = async () => {
    await saveSession();
  };

  const handleClearSession = async () => {
    if (confirm('セッションデータを削除しますか？入力内容は失われます。')) {
      await clearSession();
    }
  };

  const handleExtendSession = async () => {
    await extendSession();
  };

  const handleToggleAutoSave = () => {
    if (isAutoSaveEnabled) {
      disableAutoSave();
    } else {
      enableAutoSave();
    }
  };

  return (
    <div className={`session-status ${className}`} data-testid={testId}>
      <div className="session-status__indicator">
        <div className={`session-status__icon ${sessionExists ? 'session-status__icon--active' : 'session-status__icon--inactive'}`}>
          {isSessionLoading ? (
            <LoadingSpinner size="small" />
          ) : sessionExists ? (
            <span role="img" aria-label="セッション有効">💾</span>
          ) : (
            <span role="img" aria-label="セッション無効">📝</span>
          )}
        </div>
        
        <div className="session-status__text">
          <span className="session-status__label">
            {isSessionLoading ? '保存中...' : sessionExists ? 'セッション保存済み' : 'セッション未保存'}
          </span>
          
          {showDetails && (
            <div className="session-status__details">
              <span className="session-status__auto-save">
                自動保存: {isAutoSaveEnabled ? 'ON' : 'OFF'}
              </span>
            </div>
          )}
        </div>
      </div>

      {sessionError && (
        <div className="session-status__error" role="alert">
          <span className="session-status__error-icon">⚠️</span>
          <span className="session-status__error-text">
            セッション保存エラー: {sessionError.message}
          </span>
        </div>
      )}

      {timeoutWarning.show && (
        <div className="session-status__warning" role="alert">
          <span className="session-status__warning-icon">⏰</span>
          <span className="session-status__warning-text">
            セッション期限まで残り{timeoutWarning.remainingMinutes}分
          </span>
          <Button
            size="small"
            variant="outline"
            onClick={handleExtendSession}
            data-testid="extend-session-quick"
          >
            延長
          </Button>
        </div>
      )}

      {showDetails && (
        <div className="session-status__controls">
          <Button
            size="small"
            variant="outline"
            onClick={handleManualSave}
            disabled={isSessionLoading}
            data-testid="manual-save"
          >
            手動保存
          </Button>
          
          <Button
            size="small"
            variant="outline"
            onClick={handleToggleAutoSave}
            data-testid="toggle-auto-save"
          >
            自動保存{isAutoSaveEnabled ? 'OFF' : 'ON'}
          </Button>
          
          {sessionExists && (
            <Button
              size="small"
              variant="danger"
              onClick={handleClearSession}
              disabled={isSessionLoading}
              data-testid="clear-session"
            >
              セッション削除
            </Button>
          )}
        </div>
      )}
    </div>
  );
};

export default SessionStatus;