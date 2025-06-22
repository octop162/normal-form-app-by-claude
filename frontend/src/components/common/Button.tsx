// Reusable button component
import React from 'react';

interface ButtonProps {
  children: React.ReactNode;
  type?: 'button' | 'submit' | 'reset';
  variant?: 'primary' | 'secondary' | 'danger' | 'outline';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
  className?: string;
  'data-testid'?: string;
}

const Button: React.FC<ButtonProps> = ({
  children,
  type = 'button',
  variant = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  onClick,
  className = '',
  'data-testid': testId,
  ...props
}) => {
  const isDisabled = disabled || loading;
  
  const baseClasses = 'button';
  const variantClass = `button--${variant}`;
  const sizeClass = `button--${size}`;
  const disabledClass = isDisabled ? 'button--disabled' : '';
  const loadingClass = loading ? 'button--loading' : '';
  
  const buttonClasses = [
    baseClasses,
    variantClass,
    sizeClass,
    disabledClass,
    loadingClass,
    className
  ].filter(Boolean).join(' ');

  return (
    <button
      type={type}
      className={buttonClasses}
      disabled={isDisabled}
      onClick={onClick}
      data-testid={testId}
      aria-busy={loading}
      {...props}
    >
      {loading && (
        <span className="button__spinner" aria-hidden="true">
          ⟳
        </span>
      )}
      <span className={loading ? 'button__text--loading' : 'button__text'}>
        {children}
      </span>
    </button>
  );
};

export default Button;