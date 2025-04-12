import { browser } from "$app/environment";
import { getCookie } from "$lib/utils/getCookie";
import { writable, get } from "svelte/store";

// Define team data structure
interface TeamData {
    id: string;
    name: string;
    is_ready: boolean; // Indicates if the whole team is ready
    members: Member[]; // List of team members
}

interface Member {
    id: string;
    name: string;
    is_ready: boolean;
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

// Function to update a specific player's ready status
const setPlayerReadyState = (playerId: string, isReady: boolean) => {
    // Get the current state directly from the store
    const currentData = get(team);
    
    if (currentData) {
        // Update the specific player's ready state
        const updatedMembers = currentData.members.map(member =>
            member.id === playerId ? { ...member, is_ready: isReady } : member
        );

        // Check if all players are ready
        const isAllReady = updatedMembers.every(member => member.is_ready);

        // Update team ready status based on members' readiness
        const updatedTeam = {
            ...currentData,
            members: updatedMembers,
            is_ready: isAllReady, // The whole team is ready if all members are ready
        };
        
        // Update the store
        setTeam(updatedTeam);
    }
};

// Helper function to check if user has a team
const hasTeam = (): boolean => {
    const teamData = get(team);
    return !!teamData?.id;
};

// Function to update specific team properties
const updateTeam = (updates: Partial<TeamData>) => {
    const currentData = get(team);
    if (currentData) {
        setTeam({ ...currentData, ...updates });
    }
};

// Function to specifically update team ready status
const setTeamReady = (isReady: boolean) => {
    const currentData = get(team);
    if (currentData) {
        setTeam({ ...currentData, is_ready: isReady });
    }
};

// Function to clear team data (e.g., when user logs out)
const clearTeam = () => {
    setTeam(null);
};

export {
    team,
    setTeam,
    setPlayerReadyState,
    hasTeam,
    clearTeam,
    updateTeam,
    setTeamReady,
    TEAM_TOKEN_NAME,
    type TeamData
};