import Link from "next/link";

export default function Home() {
    return (
        <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
            {/* Navigation */}
            <nav className="fixed top-0 w-full z-50 bg-slate-900/80 backdrop-blur-md border-b border-slate-700">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex items-center justify-between h-16">
                        <div className="flex items-center gap-2">
                            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-cyan-400 rounded-lg flex items-center justify-center">
                                <span className="text-white font-bold text-sm">G</span>
                            </div>
                            <span className="text-white font-semibold text-lg">Open-Gerege</span>
                        </div>
                        <div className="flex items-center gap-4">
                            <Link
                                href="https://github.com/geregecore/open-gerege"
                                target="_blank"
                                className="text-slate-300 hover:text-white transition-colors text-sm"
                            >
                                GitHub
                            </Link>
                            <Link
                                href="/login"
                                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
                            >
                                Нэвтрэх
                            </Link>
                        </div>
                    </div>
                </div>
            </nav>

            {/* Hero Section */}
            <section className="pt-32 pb-20 px-4">
                <div className="max-w-5xl mx-auto text-center">
                    <div className="inline-flex items-center gap-2 px-3 py-1 bg-blue-500/10 border border-blue-500/20 rounded-full text-blue-400 text-sm mb-6">
                        <span className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
                        Open Source
                    </div>
                    <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-white mb-6 leading-tight">
                        Монголын хөгжүүлэгчдэд
                        <br />
                        <span className="bg-gradient-to-r from-blue-400 via-cyan-400 to-teal-400 bg-clip-text text-transparent">
                            зориулсан платформ
                        </span>
                    </h1>
                    <p className="text-lg sm:text-xl text-slate-400 max-w-2xl mx-auto mb-10">
                        Enterprise-grade backend API template болон вэб аппликейшн framework.
                        Go болон Next.js дээр суурилсан, Clean Architecture зарчмаар бүтээгдсэн.
                    </p>
                    <div className="flex flex-col sm:flex-row gap-4 justify-center">
                        <Link
                            href="/login"
                            className="px-8 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-700 hover:to-cyan-700 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40"
                        >
                            Эхлүүлэх
                        </Link>
                        <Link
                            href="/register"
                            className="px-8 py-3 bg-slate-800 hover:bg-slate-700 text-white rounded-xl font-medium transition-all border border-slate-700"
                        >
                            Бүртгүүлэх
                        </Link>
                    </div>
                </div>
            </section>

            {/* Tech Stack */}
            <section className="py-16 px-4 border-y border-slate-700/50">
                <div className="max-w-5xl mx-auto">
                    <p className="text-center text-slate-500 text-sm mb-8">ТЕХНОЛОГИЙН СТЕК</p>
                    <div className="flex flex-wrap justify-center items-center gap-8 sm:gap-12">
                        {[
                            { name: "Go", version: "1.25" },
                            { name: "Fiber", version: "v2" },
                            { name: "Next.js", version: "16" },
                            { name: "React", version: "19" },
                            { name: "PostgreSQL", version: "15" },
                            { name: "Redis", version: "7" },
                        ].map((tech) => (
                            <div key={tech.name} className="flex items-center gap-2 text-slate-400">
                                <span className="font-medium text-white">{tech.name}</span>
                                <span className="text-xs bg-slate-800 px-2 py-0.5 rounded">{tech.version}</span>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Features */}
            <section className="py-20 px-4">
                <div className="max-w-6xl mx-auto">
                    <div className="text-center mb-16">
                        <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
                            Онцлог шинж чанарууд
                        </h2>
                        <p className="text-slate-400 max-w-xl mx-auto">
                            Орчин үеийн вэб аппликейшн хөгжүүлэлтэд шаардлагатай бүх үндсэн функцууд
                        </p>
                    </div>
                    <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
                        {[
                            {
                                icon: "🔐",
                                title: "Нэвтрэлт таних",
                                description: "SSO интеграци, MFA/TOTP, Refresh token rotation",
                            },
                            {
                                icon: "👥",
                                title: "RBAC эрхийн систем",
                                description: "Role-based хандалтын удирдлага, нарийвчилсан permission",
                            },
                            {
                                icon: "🏢",
                                title: "Байгууллагын удирдлага",
                                description: "Олон түвшний бүтэц, ажилтны удирдлага",
                            },
                            {
                                icon: "📊",
                                title: "API мониторинг",
                                description: "Request logging, Prometheus metrics, Health check",
                            },
                            {
                                icon: "📝",
                                title: "Audit logging",
                                description: "Бүх үйлдлийн бүртгэл, өөрчлөлтийн түүх",
                            },
                            {
                                icon: "🚀",
                                title: "Production ready",
                                description: "Docker, CI/CD, Security headers, Rate limiting",
                            },
                        ].map((feature) => (
                            <div
                                key={feature.title}
                                className="p-6 bg-slate-800/50 border border-slate-700/50 rounded-2xl hover:bg-slate-800 hover:border-slate-600 transition-all group"
                            >
                                <div className="text-3xl mb-4">{feature.icon}</div>
                                <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-blue-400 transition-colors">
                                    {feature.title}
                                </h3>
                                <p className="text-slate-400 text-sm">{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </section>

            {/* Architecture */}
            <section className="py-20 px-4 bg-slate-800/30">
                <div className="max-w-5xl mx-auto">
                    <div className="text-center mb-12">
                        <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
                            Clean Architecture
                        </h2>
                        <p className="text-slate-400">
                            Domain-driven design, тодорхой хуваагдсан давхаргууд
                        </p>
                    </div>
                    <div className="bg-slate-900/80 border border-slate-700 rounded-2xl p-8 font-mono text-sm">
                        <pre className="text-slate-300 overflow-x-auto">
{`┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                         │
│   Router  →  Middleware  →  Handlers  →  DTOs       │
├─────────────────────────────────────────────────────┤
│                  Service Layer                       │
│         AuthService  │  UserService  │  ...         │
├─────────────────────────────────────────────────────┤
│                Repository Layer                      │
│         AuthRepo  │  UserRepo  │  OrgRepo           │
├─────────────────────────────────────────────────────┤
│                  Domain Layer                        │
│            User  │  Role  │  Organization           │
├─────────────────────────────────────────────────────┤
│                   Data Layer                         │
│              PostgreSQL  │  Redis                   │
└─────────────────────────────────────────────────────┘`}
                        </pre>
                    </div>
                </div>
            </section>

            {/* CTA */}
            <section className="py-20 px-4">
                <div className="max-w-3xl mx-auto text-center">
                    <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
                        Эхлэхэд бэлэн үү?
                    </h2>
                    <p className="text-slate-400 mb-8">
                        Бүртгүүлээд Open-Gerege платформыг туршиж үзээрэй
                    </p>
                    <div className="flex flex-col sm:flex-row gap-4 justify-center">
                        <Link
                            href="/register"
                            className="px-8 py-3 bg-gradient-to-r from-blue-600 to-cyan-600 hover:from-blue-700 hover:to-cyan-700 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/25"
                        >
                            Үнэгүй бүртгүүлэх
                        </Link>
                        <Link
                            href="https://github.com/geregecore/open-gerege"
                            target="_blank"
                            className="px-8 py-3 bg-slate-800 hover:bg-slate-700 text-white rounded-xl font-medium transition-all border border-slate-700"
                        >
                            GitHub үзэх
                        </Link>
                    </div>
                </div>
            </section>

            {/* Footer */}
            <footer className="py-8 px-4 border-t border-slate-800">
                <div className="max-w-5xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
                    <div className="flex items-center gap-2">
                        <div className="w-6 h-6 bg-gradient-to-br from-blue-500 to-cyan-400 rounded flex items-center justify-center">
                            <span className="text-white font-bold text-xs">G</span>
                        </div>
                        <span className="text-slate-400 text-sm">
                            © 2025 Gerege Core Team
                        </span>
                    </div>
                    <div className="flex items-center gap-6 text-sm text-slate-500">
                        <Link href="https://github.com/geregecore/open-gerege" target="_blank" className="hover:text-white transition-colors">
                            GitHub
                        </Link>
                        <Link href="/login" className="hover:text-white transition-colors">
                            Нэвтрэх
                        </Link>
                        <Link href="/register" className="hover:text-white transition-colors">
                            Бүртгүүлэх
                        </Link>
                    </div>
                </div>
            </footer>
        </div>
    );
}
