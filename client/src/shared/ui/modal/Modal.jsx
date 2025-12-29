import React from "react";

const Modal = ({
    isOpen,
    onClose,
    children,
    overlayClassName,
    contentClassName,
    closeButtonClassName,
    closeLabel = "Закрыть",
    closeButtonContent = "×",
    showCloseButton = true,
    contentProps = {},
    overlayProps = {},
}) => {
    if (!isOpen) {
        return null;
    }

    const { onClick: contentClick, ...restContentProps } = contentProps;
    const { onClick: overlayClick, ...restOverlayProps } = overlayProps;

    const handleOverlayClick = (event) => {
        overlayClick?.(event);
        if (!event.defaultPrevented) {
            onClose?.();
        }
    };

    const handleContentClick = (event) => {
        event.stopPropagation();
        contentClick?.(event);
    };

    return (
        <div
            className={overlayClassName}
            onClick={handleOverlayClick}
            {...restOverlayProps}
        >
            <div
                className={contentClassName}
                onClick={handleContentClick}
                role="dialog"
                aria-modal="true"
                {...restContentProps}
            >
                {showCloseButton && (
                    <button
                        type="button"
                        className={closeButtonClassName}
                        onClick={onClose}
                        aria-label={closeLabel}
                    >
                        {closeButtonContent}
                    </button>
                )}
                {children}
            </div>
        </div>
    );
};

export default Modal;
