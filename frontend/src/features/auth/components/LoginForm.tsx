"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { authApi } from "../api";
import Cookies from 'js-cookie';

export const LoginForm = () => {
    const router = useRouter();
    const [email, setEmail] = useState("admin@gerege.mn");
    const [password, setPassword] = useState("admin123");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const response = await authApi.loginLocal({ email, password });
            // The backend returns { code, message, data: { access_token: ... } }
            // The api-client usually returns 'data' if 'code' is present.
            // So 'response' might be the inner data object already.

            const token = (response as any).access_token || (response as any).data?.access_token;

            if (token) {
                // Set cookie accessible to client
                Cookies.set('session', token, { expires: 1 }); // Backend seems to like 'session' or we use Bearer
                Cookies.set('token', token, { expires: 1 });

                // Also store in localStorage for backup
                localStorage.setItem('access_token', token);

                console.log("Login successful, token saved");
                router.push("/profile");
            } else {
                throw new Error("No access token received");
            }
        } catch (err) {
            console.error(err);
            setError(err instanceof Error ? err.message : "Login failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4 w-full max-w-sm">
            <div className="space-y-2">
                <label htmlFor="email" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Email</label>
                <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="name@example.com"
                />
            </div>
            <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">Password</label>
                <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
            </div>
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Logging in..." : "Log in"}
            </Button>
        </form>
    );
};
