import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";
import { useAuth } from "../context/AuthContext";

export default function ProtectedRoute(props) {
		const navigate = useNavigate();
		const { user } = useAuth();

		createEffect(() => {
				if (!user.loading && user() === null) {
						navigate("/login", { replace: true });
				}
		});

		return <>{props.children}</>;
}