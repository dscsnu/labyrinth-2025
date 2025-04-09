export const getCookie = (cookies: string, key: string): string | null => {
    const cookie = cookies
        .split(';')
        .find(c => c.trim().startsWith(`${key}=`));
    if (!cookie) return null;
    return cookie.split('=')[1];
}