import { browser } from "$app/environment";
import { getCookie } from "$lib/utils/getCookie";
import { writable } from "svelte/store";

const TEAM_TOKEN_NAME = 'labyrinth-gdsc-team';

// Initialize from localStorage or cookies
const getInitialTeamId = (): string | null => {
    if (!browser) return null;

    const storedId = window.localStorage.getItem(TEAM_TOKEN_NAME);
    if (storedId) {
        return storedId;
    }

    return getCookie(document.cookie, TEAM_TOKEN_NAME);
};

const team = writable<string | null>(getInitialTeamId());

// Function to update team ID
const setTeam = (teamId: string | null) => {
    if (!browser) return;

    if (teamId) {
        // Store in both localStorage and cookies
        window.localStorage.setItem(TEAM_TOKEN_NAME, teamId);
        document.cookie = `${TEAM_TOKEN_NAME}=${teamId};path=/;max-age=${60 * 60 * 24 * 365}`;
    } else {
        // Clear data if teamId is null
        window.localStorage.removeItem(TEAM_TOKEN_NAME);
        document.cookie = `${TEAM_TOKEN_NAME}=null;path=/;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
    }

    team.set(teamId);
};

// Helper function to check if user has a team
const hasTeam = (): boolean => {
    const id = getInitialTeamId();
    return !!id;
};

// Function to clear team data (e.g., when user logs out)
const clearTeam = () => {
    setTeam(null);
};

export { team, setTeam, hasTeam, clearTeam, TEAM_TOKEN_NAME };