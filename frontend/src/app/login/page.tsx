import { LoginForm } from "@/features/auth/components/LoginForm";
import Link from "next/link";

export default function LoginPage() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="w-full max-w-sm space-y-8">
                <div className="flex flex-col space-y-2 text-center">
                    <h1 className="text-2xl font-semibold tracking-tight">Login</h1>
                    <p className="text-sm text-muted-foreground">Enter your credentials below</p>
                </div>
                <LoginForm />
                <div className="text-center text-sm">
                    <Link href="/" className="underline text-muted-foreground hover:text-primary">
                        Back to Home
                    </Link>
                </div>
            </div>
        </div>
    );
}
