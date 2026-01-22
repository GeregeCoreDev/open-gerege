import { RegisterForm } from "@/features/auth/components/RegisterForm";
import Link from "next/link";

export const metadata = {
    title: "Бүртгүүлэх | Open-Gerege",
    description: "Open-Gerege платформд бүртгүүлэх",
};

export default function RegisterPage() {
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
            <main className="flex-1 flex items-center justify-center px-4 py-8">
                <div className="w-full max-w-md">
                    {/* Card */}
                    <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700/50 rounded-2xl p-8 shadow-xl">
                        {/* Header */}
                        <div className="text-center mb-6">
                            <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-lg shadow-blue-500/25">
                                <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                                </svg>
                            </div>
                            <h1 className="text-2xl font-bold text-white mb-2">
                                Бүртгүүлэх
                            </h1>
                            <p className="text-slate-400 text-sm">
                                Шинэ бүртгэл үүсгэхийн тулд доорх маягтыг бөглөнө үү
                            </p>
                        </div>

                        {/* Register Form */}
                        <RegisterForm />
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
