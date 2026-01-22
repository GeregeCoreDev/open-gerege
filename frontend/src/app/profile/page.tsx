"use client";

import { useEffect } from "react";
import { useUser } from "@/features/profile/hooks/useUser";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function ProfilePage() {
    const { user, isLoading, error, fetchProfile } = useUser();

    useEffect(() => {
        fetchProfile();
    }, [fetchProfile]);

    if (isLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex min-h-screen flex-col items-center justify-center gap-4">
                <p className="text-red-500">Error: {error}</p>
                <Button onClick={() => fetchProfile()}>Retry</Button>
                <Link href="/" className="text-blue-500 hover:underline">Go Home</Link>
            </div>
        );
    }

    if (!user) {
        return (
            <div className="flex min-h-screen flex-col items-center justify-center gap-4">
                <p>No user data found. Please log in.</p>
                <Link href="/">
                    <Button>Go to Login</Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="min-h-screen p-8 max-w-4xl mx-auto">
            <div className="bg-white shadow rounded-lg p-6">
                <div className="flex items-center gap-6 mb-8">
                    {user.profile_img_url && (
                        <img
                            src={user.profile_img_url}
                            alt="Profile"
                            className="h-24 w-24 rounded-full object-cover bg-gray-100"
                        />
                    )}
                    <div>
                        <h1 className="text-2xl font-bold">{user.last_name} {user.first_name}</h1>
                        <p className="text-gray-500">{user.email}</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-1">
                        <label className="text-sm font-medium text-gray-500">Phone</label>
                        <p className="text-gray-900">{user.phone_no}</p>
                    </div>
                    <div className="space-y-1">
                        <label className="text-sm font-medium text-gray-500">Registration No</label>
                        <p className="text-gray-900">{user.reg_no}</p>
                    </div>
                    <div className="space-y-1">
                        <label className="text-sm font-medium text-gray-500">Is Organization Rep</label>
                        <p className="text-gray-900">{user.is_org ? 'Yes' : 'No'}</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
