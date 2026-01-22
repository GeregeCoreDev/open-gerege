import { LoginButton } from "@/features/auth/components/LoginButton";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm lg:flex">
        <h1 className="text-4xl font-bold mb-8">Welcome to Refactored App</h1>
        <div className="flex flex-col gap-4 items-center">
          <p className="text-lg">Please sign in to continue</p>
          <LoginButton />
          <div className="flex items-center gap-2 mt-4">
            <span className="text-sm text-gray-500">or</span>
            <a href="/login" className="text-sm text-blue-500 hover:underline">
              Login with Email
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
