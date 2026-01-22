"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { authApi } from "../api";
import { registerSchema, type RegisterFormData } from "../schemas";
import { PasswordStrengthIndicator } from "./PasswordStrengthIndicator";

export const RegisterForm = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const {
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm<RegisterFormData>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            email: "",
            password: "",
            confirmPassword: "",
            firstName: "",
            lastName: "",
            acceptTerms: false as unknown as true,
        },
    });

    const password = watch("password");

    const onSubmit = async (data: RegisterFormData) => {
        setLoading(true);
        setError(null);

        try {
            await authApi.register({
                email: data.email,
                password: data.password,
                confirmPassword: data.confirmPassword,
                firstName: data.firstName,
                lastName: data.lastName,
                acceptTerms: true, // Schema validates this is always true
            });

            setSuccess(true);
        } catch (err) {
            console.error(err);
            if (err instanceof Error) {
                // Handle specific error messages
                if (err.message.includes("email already registered")) {
                    setError("Энэ email хаяг бүртгэлтэй байна");
                } else {
                    setError(err.message);
                }
            } else {
                setError("Бүртгэл амжилтгүй боллоо. Дахин оролдоно уу.");
            }
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div
                className="space-y-4 text-center"
                role="alert"
                aria-live="polite"
            >
                <div className="p-6 bg-green-500/10 border border-green-500/30 rounded-xl">
                    <div className="w-12 h-12 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
                        <svg className="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                    </div>
                    <h2 className="text-lg font-semibold text-green-400">
                        Бүртгэл амжилттай!
                    </h2>
                    <p className="mt-2 text-sm text-slate-400">
                        Таны email хаяг руу баталгаажуулах холбоос илгээгдлээ.
                        Email-ээ шалгаж, холбоос дээр дарж бүртгэлээ баталгаажуулна уу.
                    </p>
                </div>
                <Link
                    href="/login"
                    className="inline-block text-blue-400 hover:text-blue-300 font-medium transition-colors"
                >
                    Нэвтрэх хуудас руу буцах
                </Link>
            </div>
        );
    }

    return (
        <form
            onSubmit={handleSubmit(onSubmit)}
            className="space-y-4"
            noValidate
        >
            {/* Name Row */}
            <div className="grid grid-cols-2 gap-4">
                {/* First Name */}
                <div className="space-y-2">
                    <label
                        htmlFor="firstName"
                        className="block text-sm font-medium text-slate-300"
                    >
                        Нэр
                    </label>
                    <input
                        id="firstName"
                        type="text"
                        autoComplete="given-name"
                        aria-invalid={!!errors.firstName}
                        aria-describedby={errors.firstName ? "firstName-error" : undefined}
                        {...register("firstName")}
                        className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                        placeholder="Баяр"
                    />
                    {errors.firstName && (
                        <p id="firstName-error" className="text-xs text-red-400" role="alert">
                            {errors.firstName.message}
                        </p>
                    )}
                </div>

                {/* Last Name */}
                <div className="space-y-2">
                    <label
                        htmlFor="lastName"
                        className="block text-sm font-medium text-slate-300"
                    >
                        Овог
                    </label>
                    <input
                        id="lastName"
                        type="text"
                        autoComplete="family-name"
                        aria-invalid={!!errors.lastName}
                        aria-describedby={errors.lastName ? "lastName-error" : undefined}
                        {...register("lastName")}
                        className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                        placeholder="Дорж"
                    />
                    {errors.lastName && (
                        <p id="lastName-error" className="text-xs text-red-400" role="alert">
                            {errors.lastName.message}
                        </p>
                    )}
                </div>
            </div>

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
                <label
                    htmlFor="password"
                    className="block text-sm font-medium text-slate-300"
                >
                    Нууц үг
                </label>
                <input
                    id="password"
                    type="password"
                    autoComplete="new-password"
                    aria-invalid={!!errors.password}
                    aria-describedby={errors.password ? "password-error" : "password-strength"}
                    {...register("password")}
                    className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    placeholder="••••••••"
                />
                {errors.password && (
                    <p id="password-error" className="text-sm text-red-400" role="alert">
                        {errors.password.message}
                    </p>
                )}
                <div id="password-strength">
                    <PasswordStrengthIndicator password={password} />
                </div>
            </div>

            {/* Confirm Password */}
            <div className="space-y-2">
                <label
                    htmlFor="confirmPassword"
                    className="block text-sm font-medium text-slate-300"
                >
                    Нууц үг баталгаажуулах
                </label>
                <input
                    id="confirmPassword"
                    type="password"
                    autoComplete="new-password"
                    aria-invalid={!!errors.confirmPassword}
                    aria-describedby={errors.confirmPassword ? "confirmPassword-error" : undefined}
                    {...register("confirmPassword")}
                    className="w-full h-11 px-4 rounded-xl bg-slate-700/50 border border-slate-600 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    placeholder="••••••••"
                />
                {errors.confirmPassword && (
                    <p id="confirmPassword-error" className="text-sm text-red-400" role="alert">
                        {errors.confirmPassword.message}
                    </p>
                )}
            </div>

            {/* Terms and Conditions */}
            <div className="flex items-start space-x-3">
                <input
                    id="acceptTerms"
                    type="checkbox"
                    aria-invalid={!!errors.acceptTerms}
                    aria-describedby={errors.acceptTerms ? "acceptTerms-error" : undefined}
                    {...register("acceptTerms")}
                    className="mt-1 h-4 w-4 rounded border-slate-600 bg-slate-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-0"
                />
                <div className="grid gap-1 leading-none">
                    <label
                        htmlFor="acceptTerms"
                        className="text-sm text-slate-300"
                    >
                        Үйлчилгээний нөхцөлийг зөвшөөрч байна
                    </label>
                    <p className="text-xs text-slate-500">
                        <Link href="/terms" className="text-blue-400 hover:text-blue-300 transition-colors">
                            Үйлчилгээний нөхцөл
                        </Link>
                        {" "}болон{" "}
                        <Link href="/privacy" className="text-blue-400 hover:text-blue-300 transition-colors">
                            Нууцлалын бодлого
                        </Link>
                        -ыг уншсан.
                    </p>
                </div>
            </div>
            {errors.acceptTerms && (
                <p id="acceptTerms-error" className="text-sm text-red-400" role="alert">
                    {errors.acceptTerms.message}
                </p>
            )}

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
                        Бүртгэж байна...
                    </>
                ) : (
                    "Бүртгүүлэх"
                )}
            </button>

            {/* Login link */}
            <p className="text-center text-sm text-slate-400">
                Бүртгэлтэй юу?{" "}
                <Link
                    href="/login"
                    className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
                >
                    Нэвтрэх
                </Link>
            </p>
        </form>
    );
};
