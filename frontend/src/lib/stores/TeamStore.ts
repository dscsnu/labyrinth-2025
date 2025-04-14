import { browser } from "$app/environment";
import { getCookie } from "$lib/utils/getCookie";
import { writable, get } from "svelte/store";

// Define team data structure
export interface ITeamData {
    id: string;
    name: string;
    allReady: boolean;
    members: IMember[];
}

export interface IMember {
    id: string;
    name: string;
    isReady: boolean;
}

export const TOKEN_NAME = 'labyrinth-gdsc-team';

// Initialize from localStorage or cookies
const getInitialTeamData = (): ITeamData | null => {
    if (!browser) return null;

    // Try local storage first
    const storedData = window.localStorage.getItem(TOKEN_NAME);
    if (storedData) {
        try {
            return JSON.parse(storedData);
        } catch {}
    }

    // Try cookie as fallback
    const cookieData = getCookie(document.cookie, TOKEN_NAME);
    if (cookieData) {
        try {
            return JSON.parse(cookieData);
        } catch {}
    }

    return null;
};

export const TeamStore = writable<ITeamData | null>(getInitialTeamData());

// Function to update team data
export const setTeam = (teamData: ITeamData | null) => {
    if (!browser) return;

    if (teamData) {
        // Serialize team data to JSON
        const serialized = JSON.stringify(teamData);

        // Store in both localStorage and cookies
        window.localStorage.setItem(TOKEN_NAME, serialized);
        document.cookie = `${TOKEN_NAME}=${encodeURIComponent(serialized)};path=/;max-age=${60 * 60 * 24 * 365}`;
    } else {
        // Clear data if teamData is null
        window.localStorage.removeItem(TOKEN_NAME);
        document.cookie = `${TOKEN_NAME}=null;path=/;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
    }

    TeamStore.set(teamData);
};

// Function to update a specific player's ready status
export const setPlayerReadyState = (playerId: string, isReady: boolean) => {
    TeamStore.update(current => {
        if (!current) return current;

        const members = current.members.map(member =>
            member.id === playerId ? { ...member, isReady } : member
        );
        const allReady = members.every(m => m.isReady);

        const updated = { ...current, members, allReady };
        setTeam(updated);
        return updated;
    });
};

// Helper function to check if user has a team
export const hasTeam = (): boolean => {
    const teamData = get(TeamStore);
    return !!teamData?.id;
};

// Function to update specific team properties
export const updateTeam = (updates: Partial<ITeamData>) => {
    const currentData = get(TeamStore);
    if (currentData) {
        setTeam({ ...currentData, ...updates });
    }
};

// Function to specifically update team ready status
export const setTeamReady = (isReady: boolean) => {
    const currentData = get(TeamStore);
    if (currentData) {
        setTeam({ ...currentData, allReady: isReady });
    }
};

// Function to clear team data (e.g., when user logs out)
export const clearTeam = () => {
    setTeam(null);
};