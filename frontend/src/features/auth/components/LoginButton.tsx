"use client";

import { Button } from "@/components/ui/button";

export const LoginButton = () => {
    const handleLogin = () => {
        // Redirect directly to the backend login endpoint
        window.location.href = "/api/v1/auth/login";
    };

    return (
        <Button onClick={handleLogin}>
            Log in with SSO
        </Button>
    );
};
