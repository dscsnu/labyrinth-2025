import { customAlphabet, urlAlphabet } from "nanoid";
const nanoid = customAlphabet(urlAlphabet, 8);
export default nanoid;