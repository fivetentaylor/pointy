import axios from "axios";

import { APP_HOST } from "./urls";

export default axios.create({
  baseURL: APP_HOST,
  withCredentials: true,
});
