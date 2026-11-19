import { handleResponse } from ".";

const baseUrl = "/api/auth";

export const authService = {
	login: (data) =>
		fetch(`${baseUrl}/login`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data),
		}).then(handleResponse),

	logout: () =>
		fetch(`${baseUrl}/logout`, {
			method: "POST",
		}).then(handleResponse),
	getSelf: () => fetch(`${baseUrl}/self`).then(handleResponse),

	/* email */
	forgotPassword: (data) =>
		fetch(`${baseUrl}/forgot`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data),
		}).then(handleResponse),

	/* token, newPassword */ /* tokenin saa ?token=... */
	resetPassword: (data) =>
		fetch(`${baseUrl}/reset`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data),
		}).then(handleResponse),
};
