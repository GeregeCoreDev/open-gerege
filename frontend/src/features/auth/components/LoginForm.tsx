"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { authApi } from "../api";
import { loginSchema, type LoginFormData } from "../schemas";
import type { LoginResponse } from "../types";
import Cookies from "js-cookie";

export const LoginForm = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Get redirect URL from query params
    const redirectUrl = searchParams.get("redirect") || "/profile";

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            email: "",
            password: "",
            rememberMe: false,
        },
    });

    const onSubmit = async (data: LoginFormData) => {
        setLoading(true);
        setError(null);

        try {
            const response = await authApi.loginLocal({
                email: data.email,
                password: data.password,
            });

            // The backend returns { code, message, data: { access_token: ... } }
            // The api-client usually returns 'data' if 'code' is present.
            // Extract token from response
            const responseData = response as LoginResponse & { data?: LoginResponse };
            const token = responseData.access_token || responseData.data?.access_token;

            if (token) {
                // Set cookie expiry based on rememberMe
                const expires = data.rememberMe ? 7 : 1; // 7 days or 1 day

                // Set cookie accessible to client
                Cookies.set("session", token, { expires });
                Cookies.set("token", token, { expires });

                // Also store in localStorage for backup
                localStorage.setItem("access_token", token);

                console.log("Login successful, token saved");
                router.push(redirectUrl);
            } else {
                throw new Error("No access token received");
            }
        } catch (err) {
            console.error(err);
            if (err instanceof Error) {
                // Handle specific error messages
                if (err.message.includes("invalid")) {
                    setError("Email эсвэл нууц үг буруу байна");
                } else if (err.message.includes("locked")) {
                    setError("Таны бүртгэл түр түгжигдсэн байна. Дараа дахин оролдоно уу.");
                } else if (err.message.includes("not active")) {
                    setError("Таны бүртгэл идэвхгүй байна. Админтай холбогдоно уу.");
                } else {
                    setError(err.message);
                }
            } else {
                setError("Нэвтрэлт амжилтгүй боллоо");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <form
            onSubmit={handleSubmit(onSubmit)}
            className="space-y-5"
            noValidate
        >
            {/* Email */}
            <div className="space-y-2">
                <label
                    htmlFor="email"
                    className="block text-sm font-medium text-slate-300"
                >
                    Email хаяг
                </label>
                <input
                    id="email"
                    type="email"
                    autoComplete="email"
                    aria-invalid={!!errors.email}
                    aria-describedby={errors.email ? "email-error" : undefined}
                    {...register("email")}
                    className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    placeholder="name@example.com"
                />
                {errors.email && (
                    <p id="email-error" className="text-sm text-red-400" role="alert">
                        {errors.email.message}
                    </p>
                )}
            </div>

            {/* Password */}
            <div className="space-y-2">
                <div className="flex items-center justify-between">
                    <label
                        htmlFor="password"
                        className="block text-sm font-medium text-slate-300"
                    >
                        Нууц үг
                    </label>
                    <Link
                        href="/forgot-password"
                        className="text-sm text-blue-400 hover:text-blue-300 transition-colors"
                    >
                        Нууц үг мартсан?
                    </Link>
                </div>
                <input
                    id="password"
                    type="password"
                    autoComplete="current-password"
                    aria-invalid={!!errors.password}
                    aria-describedby={errors.password ? "password-error" : undefined}
                    {...register("password")}
                    className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    placeholder="••••••••"
                />
                {errors.password && (
                    <p id="password-error" className="text-sm text-red-400" role="alert">
                        {errors.password.message}
                    </p>
                )}
            </div>

            {/* Remember Me */}
            <div className="flex items-center">
                <input
                    id="rememberMe"
                    type="checkbox"
                    {...register("rememberMe")}
                    className="h-4 w-4 rounded border-slate-600 bg-slate-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-0"
                />
                <label
                    htmlFor="rememberMe"
                    className="ml-2 text-sm text-slate-400"
                >
                    Намайг сана
                </label>
            </div>

            {/* Error message */}
            {error && (
                <div
                    className="p-4 rounded-xl bg-red-500/10 border border-red-500/30"
                    role="alert"
                >
                    <p className="text-sm text-red-400 flex items-center gap-2">
                        <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {error}
                    </p>
                </div>
            )}

            {/* Submit button */}
            <button
                type="submit"
                disabled={loading}
                className="w-full h-11 px-4 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-700 hover:to-cyan-700 disabled:from-blue-600/50 disabled:to-cyan-600/50 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 disabled:shadow-none flex items-center justify-center gap-2"
            >
                {loading ? (
                    <>
                        <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Нэвтэрч байна...
                    </>
                ) : (
                    "Нэвтрэх"
                )}
            </button>
        </form>
    );
};
