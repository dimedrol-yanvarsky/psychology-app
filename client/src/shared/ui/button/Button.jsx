import React from "react";
import { Link } from "react-router-dom";

const Button = ({ as = "button", to, className, children, ...props }) => {
    if (as === "link") {
        return (
            <Link to={to} className={className} {...props}>
                {children}
            </Link>
        );
    }

    return (
        <button className={className} {...props}>
            {children}
        </button>
    );
};

export default Button;
