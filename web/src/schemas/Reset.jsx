import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import { useNavigate } from "@solidjs/router";
import * as v from "valibot";

const resetSchema = v.pipe(
		v.object({
				newPassword: v.pipe(
						v.string(),
						v.minLength(3, "Salasanassa tulee olla vähintään 3 merkkiä"),
						v.maxLength(50, "Salasanassa voi olla enintään 50 merkkiä"),
				),
		})
);

function ResetForm(props) {
		const navigate = useNavigate();
		
		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const resetForm = createForm({
				schema: resetSchema,
		});

		return (
				<Form
						of={resetForm}
						onSubmit={handleSubmit}
						class="form"
				>
						<div style="padding: 15px;">
								<Field
										of={resetForm}
										path={["newPassword"]}
								>
										{(field) => (
												<div>
														<label for={field.props.name}>Uusi salasana:</label>
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
						</div>
						<div style="padding: 15px; border-top: 1px solid var(--border);">
								<button
										class="btn-primary"
										type="submit"
								>
										Nollaa salasana
								</button>
						</div>
				</Form>
		);
}

export default ResetForm;