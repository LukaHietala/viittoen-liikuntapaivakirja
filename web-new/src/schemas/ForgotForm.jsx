import { createForm, Field, Form, reset, setInput, getInput } from "@formisch/solid";
import * as v from "valibot";

const forgotSchema = v.pipe(
		v.object({
				email: v.pipe(
						v.string(),
						v.minLength(3, "Sähköpostissa tulee olla vähintään 3 merkkiä"),
						v.email(),
						v.maxLength(200, "Sähköpostissa voi olla enintään 200 merkkiä"),
				),
		})
);

function ForgotForm(props) {
		const handleSubmit = async (data) => {
				await props.onSubmit(data);
		};

		const forgotForm = createForm({
				schema: forgotSchema,
		});

		return (
				<Form
						of={forgotForm}
						onSubmit={handleSubmit}
						class="form"
				>				
						<Field
								of={forgotForm}
								path={["email"]}
						>
								{(field) => (
										<div>
												<label for={field.props.name}>Sähköposti:</label>
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
						
						<button
								class="btn-primary"
								type="submit"
						>
								Lähetä sähköposti
						</button>
				</Form>
		);
}

export default ForgotForm;