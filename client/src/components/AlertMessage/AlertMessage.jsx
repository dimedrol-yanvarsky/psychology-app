import React, { useState, useRef } from "react";
import styles from "./AlertMessage.module.css";

    const AlertMessage = (props) => {
        return (
            <div className={`${styles.alert} ${props.status === "error" ? styles.error : styles.success}`}>
                <p className={styles.message}>{props.message}</p>
            </div>
        );
    };

export default AlertMessage;
