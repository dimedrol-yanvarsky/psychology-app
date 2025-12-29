import React from "react";
import styles from "./AlertMessage.module.css";

const AlertMessage = (props) => {
    return (
        <div
            className={`${styles.alert} ${
                props.statusAlert === "error" ? styles.error : styles.success
            }`}
        >
            <p className={styles.message}>{props.messageAlert}</p>
        </div>
    );
};

export default AlertMessage;
