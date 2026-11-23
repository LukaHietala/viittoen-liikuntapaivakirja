import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";
import { useAuth } from "../context/AuthContext";

export default function AdminRoute(props) {
		const navigate = useNavigate();
		const { user } = useAuth();

		createEffect(() => {
				if (!user.loading && user() && !user().isAdmin) {
						navigate("/", { replace: true });
				}
		});

		return <>{props.children}</>;
}