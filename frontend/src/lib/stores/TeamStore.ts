import createPersistentStore from "$lib/utils/createPersistentStore";
import { get } from "svelte/store";

export interface ITeamData {
    id: string;
    name: string;
    members: IMember[];
}

export interface IMember {
    id: string;
    name: string;
    email: string;
    isReady: boolean;
}

export const TOKEN_NAME = 'labyrinth-gdsc-team';

const {
    store: TeamStore, set: setTeam
} = createPersistentStore<ITeamData>(TOKEN_NAME);

export const setPlayerReadyState = (playerId: string, isReady: boolean) => {
    const current = get(TeamStore);
    if (!current) return;

    const members = current.members.map(member =>
        member.id === playerId ? { ...member, isReady } : member
    );

    setTeam({ ...current, members });
}

export const updateTeam = (update: Partial<ITeamData>) => {
    const current = get(TeamStore);
    if (!current) return;

    setTeam({ ...current, ...update });
}

export const clearTeam = () => setTeam(null);

export { TeamStore, setTeam };