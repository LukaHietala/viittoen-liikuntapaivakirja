import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import { useNavigate } from "@solidjs/router";
import * as v from "valibot";

const loginSchema = v.pipe(
		v.object({
				name: v.pipe(
						v.string(),
						v.minLength(3, "Nimessä tulee olla vähintään 3 merkkiä"),
						v.maxLength(50, "Nimessä voi olla enintään 50 merkkiä"),
				),
				password: v.pipe(
						v.string(),
						v.minLength(3, "Salasanassa tulee olla vähintään 3 merkkiä"),
						v.maxLength(50, "Salasanassa voi olla enintään 50 merkkiä"),
				),
		})
);

function LoginForm(props) {
		const navigate = useNavigate();
		
		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const loginForm = createForm({
				schema: loginSchema,
		});

		return (
				<Form
						of={loginForm}
						onSubmit={handleSubmit}
						class="form"
				>				
						<Field
								of={loginForm}
								path={["name"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Nimi:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="text"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>

						<Field
								of={loginForm}
								path={["password"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Salasana:</label>
												<br />
												<input
														{...field.props}
														id={field.props.name}
														value={field.input}
														type="password"
														required
												/>
												{field.errors && <small class="form-error">{field.errors[0]}</small>}
										</div>
								)}
						</Field>
						
						<button
								class="btn-primary"
								type="submit"
						>
								Kirjaudu
						</button>
						<button
								onClick={() => navigate("/forgot")}
								type="button" >
								unohditko salasani?
						</button>
				</Form>
		);
}

export default LoginForm;