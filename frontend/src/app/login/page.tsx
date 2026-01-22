import { LoginForm } from "@/features/auth/components/LoginForm";
import Link from "next/link";

export const metadata = {
    title: "Нэвтрэх | Open-Gerege",
    description: "Open-Gerege платформд нэвтрэх",
};

export default function LoginPage() {
    return (
        <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 flex flex-col">
            {/* Header */}
            <header className="p-4">
                <Link href="/" className="flex items-center gap-2 w-fit">
                    <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-lg flex items-center justify-center">
                        <span className="text-white font-bold text-sm">G</span>
                    </div>
                    <span className="text-white font-semibold text-lg">Open-Gerege</span>
                </Link>
            </header>

            {/* Main Content */}
            <main className="flex-1 flex items-center justify-center px-4 py-12">
                <div className="w-full max-w-md">
                    {/* Card */}
                    <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700/50 rounded-2xl p-8 shadow-xl">
                        {/* Header */}
                        <div className="text-center mb-8">
                            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-blue-500/25">
                                <span className="text-white font-bold text-2xl">G</span>
                            </div>
                            <h1 className="text-2xl font-bold text-white mb-2">
                                Тавтай морил
                            </h1>
                            <p className="text-slate-400 text-sm">
                                Open-Gerege платформд нэвтрэхийн тулд мэдээллээ оруулна уу
                            </p>
                        </div>

                        {/* Login Form */}
                        <LoginForm />

                        {/* Divider */}
                        <div className="relative my-6">
                            <div className="absolute inset-0 flex items-center">
                                <div className="w-full border-t border-slate-700"></div>
                            </div>
                            <div className="relative flex justify-center text-xs uppercase">
                                <span className="bg-slate-800/50 px-2 text-slate-500">
                                    эсвэл
                                </span>
                            </div>
                        </div>

                        {/* SSO Button */}
                        <button
                            type="button"
                            className="w-full px-4 py-3 bg-slate-700/50 hover:bg-slate-700 text-white rounded-xl font-medium transition-all border border-slate-600 flex items-center justify-center gap-2"
                        >
                            <svg
                                className="w-5 h-5"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                                />
                            </svg>
                            Gerege SSO-р нэвтрэх
                        </button>

                        {/* Register Link */}
                        <p className="text-center text-sm text-slate-400 mt-6">
                            Бүртгэлгүй юу?{" "}
                            <Link
                                href="/register"
                                className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
                            >
                                Бүртгүүлэх
                            </Link>
                        </p>
                    </div>

                    {/* Back to Home */}
                    <div className="text-center mt-6">
                        <Link
                            href="/"
                            className="text-slate-500 hover:text-slate-300 text-sm transition-colors inline-flex items-center gap-1"
                        >
                            <svg
                                className="w-4 h-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M10 19l-7-7m0 0l7-7m-7 7h18"
                                />
                            </svg>
                            Нүүр хуудас руу буцах
                        </Link>
                    </div>
                </div>
            </main>

            {/* Footer */}
            <footer className="p-4 text-center text-slate-600 text-sm">
                © 2025 Gerege Core Team
            </footer>
        </div>
    );
}
