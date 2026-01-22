import { create } from 'zustand';
import { profileApi } from '../api';
import { UserProfile } from '../types';

interface UserStore {
    user: UserProfile | null;
    isLoading: boolean;
    error: string | null;
    fetchProfile: () => Promise<void>;
    logout: () => void;
}

export const useUser = create<UserStore>((set) => ({
    user: null,
    isLoading: false,
    error: null,
    fetchProfile: async () => {
        set({ isLoading: true, error: null });
        try {
            const user = await profileApi.getProfile();
            set({ user, isLoading: false });
        } catch (error) {
            set({
                error: error instanceof Error ? error.message : 'Failed to fetch profile',
                isLoading: false,
                user: null
            });
        }
    },
    logout: () => set({ user: null, error: null }),
}));
