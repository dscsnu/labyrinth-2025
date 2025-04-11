import { browser } from "$app/environment";
import { getCookie } from "$lib/utils/getCookie";
import { writable } from "svelte/store";

// Define team data structure
interface TeamData {
    id: string;
    name: string;
    isReady: boolean;
}

const TEAM_TOKEN_NAME = 'labyrinth-gdsc-team';

// Initialize from localStorage or cookies
const getInitialTeamData = (): TeamData | null => {
    if (!browser) return null;

    // Try local storage first
    const storedData = window.localStorage.getItem(TEAM_TOKEN_NAME);
    if (storedData) {
        try {
            return JSON.parse(storedData);
        } catch (e) {
            // If JSON parsing fails, continue to cookie fallback
        }
    }

    // Try cookie as fallback
    const cookieData = getCookie(document.cookie, TEAM_TOKEN_NAME);
    if (cookieData) {
        try {
            return JSON.parse(cookieData);
        } catch (e) {
            // If JSON parsing fails, return null
        }
    }

    return null;
};

const team = writable<TeamData | null>(getInitialTeamData());

// Function to update team data
const setTeam = (teamData: TeamData | null) => {
    if (!browser) return;

    if (teamData) {
        // Serialize team data to JSON
        const serialized = JSON.stringify(teamData);

        // Store in both localStorage and cookies
        window.localStorage.setItem(TEAM_TOKEN_NAME, serialized);
        document.cookie = `${TEAM_TOKEN_NAME}=${encodeURIComponent(serialized)};path=/;max-age=${60 * 60 * 24 * 365}`;
    } else {
        // Clear data if teamData is null
        window.localStorage.removeItem(TEAM_TOKEN_NAME);
        document.cookie = `${TEAM_TOKEN_NAME}=null;path=/;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
    }

    team.set(teamData);
};

// Helper function to check if user has a team
const hasTeam = (): boolean => {
    const teamData = getInitialTeamData();
    return !!teamData?.id;
};

// Function to update specific team properties
const updateTeam = (updates: Partial<TeamData>) => {
    const currentData = getInitialTeamData();
    if (currentData) {
        setTeam({ ...currentData, ...updates });
    }
};

// Function to specifically update team ready status
const setTeamReady = (isReady: boolean) => {
    const currentData = getInitialTeamData();
    if (currentData) {
        setTeam({ ...currentData, isReady });
    }
};

// Function to clear team data (e.g., when user logs out)
const clearTeam = () => {
    setTeam(null);
};

export {
    team,
    setTeam,
    hasTeam,
    clearTeam,
    updateTeam,
    setTeamReady,
    TEAM_TOKEN_NAME,
    type TeamData
};