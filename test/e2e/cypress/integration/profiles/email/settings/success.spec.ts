// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

import { appPrefix, gen, website } from "../../../../helpers"
import { routes as express } from "../../../../helpers/express"
import { routes as react } from "../../../../helpers/react"

context("Settings success with email profile", () => {
  ;[
    {
      route: express.settings,
      base: express.base,
      app: "express" as "express",
      profile: "email",
      login: express.login,
    },
    {
      route: react.settings,
      base: react.base,
      app: "react" as "react",
      profile: "spa",
      login: react.login,
    },
  ].forEach(({ route, profile, app, base, login }) => {
    describe(`for app ${app}`, () => {
      let email = gen.email()
      const firstPassword = gen.password()
      const secondPassword = gen.password()
      const thirdPassword = gen.password()
      let password = firstPassword

      const up = (value: string) => `not-${value}`
      const down = (value: string) => value.replace(/not-/, "")

      before(() => {
        cy.useConfigProfile(profile)
        cy.registerApi({
          email,
          password,
          fields: { "traits.website": website },
        })
        cy.proxy(app)
      })

      beforeEach(() => {
        cy.deleteMail()
        cy.login({ email, password, cookieUrl: base })
        cy.visit(route)
      })

      it("shows all settings forms", () => {
        cy.get(appPrefix(app) + "h3").should("contain.text", "Profile")
        cy.get('input[name="traits.email"]').should("contain.value", email)
        cy.get('input[name="traits.website"]').should("contain.value", website)

        cy.get("h3").should("contain.text", "Password")
        cy.get('input[name="password"]').should("be.empty")
      })

      describe("password", () => {
        it("modifies the password with privileged session", () => {
          // Once input weak password to test which error message is cleared after updating successfully
          cy.get('input[name="password"]').clear().type("123")
          cy.get('button[value="password"]').click()
          cy.get('[data-testid="ui/message/1050001"]').should("not.exist")
          cy.get('[data-testid="ui/message/4000032"]').should("exist")
          cy.get('input[name="password"]').should("be.empty")

          password = secondPassword
          cy.get('input[name="password"]').clear().type(secondPassword)
          cy.get('button[value="password"]').click()
          cy.expectSettingsSaved()
          cy.get('[data-testid="ui/message/4000032"]').should("not.exist")
          cy.get('[data-testid="ui/message/1050001"]').should("exist")
          cy.get('input[name="password"]').should("be.empty")
        })

        it("is unable to log in with the old password", () => {
          cy.login({
            email: email,
            password: firstPassword,
            expectSession: false,
            cookieUrl: base,
          })
        })

        it("modifies the password with an unprivileged session", () => {
          password = thirdPassword
          cy.get('input[name="password"]').clear().type(password)
          cy.shortPrivilegedSessionTime() // wait for the privileged session to time out
          cy.get('button[value="password"]').click()

          cy.reauth({ expect: { email }, type: { password: secondPassword } })

          cy.url().should("include", "/settings")
          cy.expectSettingsSaved()
          cy.get('input[name="password"]').should("be.empty")
        })
      })

      describe("profile", () => {
        it("modifies an unprotected traits", () => {
          cy.get('input[name="traits.website"]')
            .clear()
            .type("https://github.com/ory")
          cy.get('input[name="traits.age"]').clear().type("30")
          cy.get('input[type="checkbox"][name="traits.tos"]').click({
            force: true,
          })
          cy.submitProfileForm()
          cy.expectSettingsSaved()

          cy.get('input[name="traits.website"]').should(
            "contain.value",
            "https://github.com/ory",
          )
          cy.get('input[type="checkbox"][name="traits.tos"]')
            .should("be.checked")
            .click({ force: true })
          cy.get('input[name="traits.age"]')
            .should("have.value", "30")
            .clear()
            .type("90")

          cy.submitProfileForm()
          cy.expectSettingsSaved()

          cy.get('input[type="checkbox"][name="traits.tos"]').should(
            "not.be.checked",
          )
          cy.get('input[name="traits.age"]').should("have.value", "90")
        })

        it("modifies a protected trait with privileged session", () => {
          email = up(email)
          cy.disableVerification()
          cy.get('input[name="traits.email"]').clear().type(email)
          cy.get('button[value="profile"]').click()
          cy.expectSettingsSaved()
          cy.get('input[name="traits.email"]').should("contain.value", email)
        })

        it("is unable to log in with the old email", () => {
          cy.login({
            email: down(email),
            password,
            expectSession: false,
            cookieUrl: base,
          })
        })

        it("modifies a protected trait with unprivileged session", () => {
          email = up(email)
          cy.get('input[name="traits.email"]').clear().type(email)
          cy.shortPrivilegedSessionTime() // wait for the privileged session to time out
          cy.get('button[value="profile"]').click()

          cy.reauth({ expect: { email: down(email) }, type: { password } })

          cy.url().should("include", "/settings")
          cy.expectSettingsSaved()
          cy.get('input[name="traits.email"]').should("contain.value", email)
        })

        if (app === "react") {
          it("shows verification screen after email update", () => {
            cy.deleteMail()
            cy.enableVerification()
            email = up(email)
            cy.get('input[name="traits.email"]').clear().type(email)
            cy.get('button[value="profile"]').click()

            cy.url().should("contain", "verification")
            cy.getVerificationCodeFromEmail(email).then((code) => {
              cy.get("input[name=code]").type(code)
              cy.get("button[name=method][value=code]").click()
            })

            cy.get('[data-testid="ui/message/1080002"]').should(
              "have.text",
              "You successfully verified your email address.",
            )
          })
        }
      })
    })
  })
})
