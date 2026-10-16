import { error, json } from '@sveltejs/kit';
import jwt from 'jsonwebtoken';
import { JWT_SESSION_KEY } from "$env/static/private";

export async function POST({ request }) {
    try {
        const data = await request.json();
        if (!data.payload) {
            throw error(400, 'Payload is required');
        }

        const currentTime = Math.floor(Date.now() / 1000);

        const tokenData = {
            payload: data.payload,
            timestamp: currentTime
        };


        const token = jwt.sign(tokenData, JWT_SESSION_KEY, {
            expiresIn: '1h'
        });

        return json({
            success: true,
            token: token
        });
    } catch (err) {
        console.error('Error encrypting payload:', err);
        return error(500, 'Failed to encrypt payload');
    }
}