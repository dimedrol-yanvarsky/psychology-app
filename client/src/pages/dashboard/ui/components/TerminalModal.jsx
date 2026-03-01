import React from "react";
import { Terminal } from "../../../../widgets/terminal";
import { Modal } from "../../../../shared/ui/modal";
import styles from "../DashboardPage.module.css";

const TerminalModal = ({ isOpen, onClose, setIsTerminalOpen }) => {
    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            overlayClassName={styles.terminalOverlay}
            contentClassName={styles.terminalModal}
            showCloseButton={false}
        >
            <Terminal setIsTerminalOpen={setIsTerminalOpen} />
        </Modal>
    );
};

export default TerminalModal;
