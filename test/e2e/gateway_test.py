import requests
import unittest
import common

"""
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
"""

GATEWAY_URL = "http://localhost/v0"

class TestAccount(unittest.TestCase):

    def test_signup(self):
        print("Testing signup")
        common.make_user(GATEWAY_URL, common.generate_user())

    def test_login(self):
        print("Testing login")
        user = common.generate_user()
        common.make_user(GATEWAY_URL, user)

        creds = {"email": user["email"],"password": user["password"]}
        common.login(GATEWAY_URL, creds)

    def test_logout(self):
        print("Testing logout")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        common.logout(GATEWAY_URL, auth)

    def test_getsession(self):
        print("Testing get session")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        common.getsession(GATEWAY_URL, auth)

    
class TestForm(unittest.TestCase):

    def test_make_form(self):
        print("Testing make form")

        auth = common.make_user(GATEWAY_URL, common.generate_user())
        common.make_form(GATEWAY_URL, auth, common.generate_form())

    def test_get_forms(self):
        print("Testing get forms")

        auth = common.make_user(GATEWAY_URL, common.generate_user())
        common.make_form(GATEWAY_URL, auth, common.generate_form())
        common.make_form(GATEWAY_URL, auth, common.generate_form())
        
        assert len(common.get_forms(GATEWAY_URL, auth)) == 2

    def test_get_form(self):
        print("Testing get form")

        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.get_form(GATEWAY_URL, auth, form["id"])

    def test_get_form_user(self):
        print("Testing get form: user")

        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.get_form_client(GATEWAY_URL, form["id"])

    def test_post_form_user(self):
        print("Testing post form: user")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())

    def test_get_response(self):
        print("Testing get response")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        resps = common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})
        
        common.get_response(GATEWAY_URL, auth, resps[0]["id"])

    def test_get_responses(self):
        print("Testing get responses (no filter)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        common.post_form_user(GATEWAY_URL, form2["id"], common.generate_response())
        common.post_form_user(GATEWAY_URL, form2["id"], common.generate_response())

        assert len(common.get_responses(GATEWAY_URL, auth)) == 5

    def test_get_responses_subject(self):
        print("Testing get responses (filter: subject)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()


        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"subject": resp1["subject"]})) == 3

    def test_get_responses_email(self):
        print("Testing get responses (filter: emailAddress)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"emailAddress": resp1["email"]})) == 3

    def test_get_responses_active(self):
        print("Testing get responses (filter: activeOnly)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        # mark 3 of the responses as non-active, leaving 2 active
        resps = common.get_responses(GATEWAY_URL, auth)
        for resp in resps[:3]:
            common.patch_response(GATEWAY_URL, auth, resp["id"], False)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"activeOnly": True})) == 2

    def test_get_responses_formID(self):
        print("Testing get responses (filter: formID)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})) == 3

    def test_get_responses_name(self):
        print("Testing get responses (filter: name)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"name": resp1["name"]})) == 3

    def test_get_responses_today(self):
        print("Testing get responses (filter: todayOnly)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"todayOnly": True})) == 5

    def test_get_responses_inquiry(self):
        print("Testing get responses (filter: inquiryType)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp1["inquiryType"] = "casework"

        resp2 = common.generate_response()
        resp2["inquiryType"] = "general"

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"inquiryType": "casework"})) == 3
        assert len(common.get_responses_params(GATEWAY_URL, auth, {"inquiryType": "general"})) == 2
        assert len(common.get_responses_params(GATEWAY_URL, auth, {"inquiryType": "INVALID"})) == 0

    def test_get_responses_multi(self):
        print("Testing get responses (filter: multiple)")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        resp1 = common.generate_response()
        resp2 = common.generate_response()

        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp1)
        common.post_form_user(GATEWAY_URL, form["id"], resp2)
        common.post_form_user(GATEWAY_URL, form2["id"], resp1)
        common.post_form_user(GATEWAY_URL, form2["id"], resp2)

        assert len(common.get_responses_params(GATEWAY_URL, auth, {"name": resp1["name"], "formID": form["id"]})) == 2

    def test_patch_response(self):
        print("Testing patch response")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        resp = common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})[0]

        common.patch_response(GATEWAY_URL, auth, resp["id"], False)

        updated = common.get_response(GATEWAY_URL, auth, resp["id"])
        assert updated["active"] == False

        common.patch_response(GATEWAY_URL, auth, resp["id"], True)
        updated = common.get_response(GATEWAY_URL, auth, resp["id"])
        assert updated["active"] == True

    def test_create_tag(self):
        print("Testing create tag")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        resp = common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})[0]

        common.make_tag(GATEWAY_URL, auth, resp["id"], common.randstr(5, 10))

    def test_delete_tag(self):
        print("Testing delete tag")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        resp = common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})[0]

        tag_val = common.randstr(5, 10)
        common.make_tag(GATEWAY_URL, auth, resp["id"], tag_val)
        common.delete_tag(GATEWAY_URL, auth, resp["id"], tag_val)

    def test_get_all_tags(self):
        print("Testing get tags")

        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form = common.make_form(GATEWAY_URL, auth, common.generate_form())
        common.post_form_user(GATEWAY_URL, form["id"], common.generate_response())
        resp = common.get_responses_params(GATEWAY_URL, auth, {"formID": form["id"]})[0]

      

        assert len(common.get_all_tags(GATEWAY_URL, auth)) == 3

    def test_get_tags_response(self):
        print("Testing get tags by response id")
        auth = common.make_user(GATEWAY_URL, common.generate_user())
        form1 = common.make_form(GATEWAY_URL, auth, common.generate_form())
        form2 = common.make_form(GATEWAY_URL, auth, common.generate_form())

        common.post_form_user(GATEWAY_URL, form1["id"], common.generate_response())
        common.post_form_user(GATEWAY_URL, form2["id"], common.generate_response())

        resp1 = common.get_responses_params(GATEWAY_URL, auth, {"formID": form1["id"]})[0]
        resp2 = common.get_responses_params(GATEWAY_URL, auth, {"formID": form2["id"]})[0]


        assert len(common.get_all_tags_response(GATEWAY_URL, auth, resp1["id"])) == 3
        assert len(common.get_all_tags_response(GATEWAY_URL, auth, resp2["id"])) == 3
        assert len(common.get_response(GATEWAY_URL, auth, resp1["id"])["tags"]) == 3
        assert len(common.get_response(GATEWAY_URL, auth, resp2["id"])["tags"]) == 3

        allresps = common.get_responses(GATEWAY_URL, auth)
        assert sum([len(res["tags"]) for res in allresps]) == 6


class TestMailto(unittest.TestCase):

    def test_simple_to(self):
        print("Testing mailto simple")

        body = {
            "to": ["test@example.com"]
        }

        assert common.post_mailto(GATEWAY_URL, body) == r'mailto:test@example.com'

    def test_multiple_to(self):
        print("Testing mailto multiple to's")
        body = {
            "to": ["test@example.com", "test1@example.com", "test2@example.com"]
        }

        assert common.post_mailto(GATEWAY_URL, body) == r'mailto:test@example.com,test1@example.com,test2@example.com'

    def test_to_subject_body(self):
        print("Testing mailto (to, subject, body)")
        body = {
            "to": ["test@example.com"],
            "subject": "test subject",
            "body": "test body"
        }

        assert common.post_mailto(GATEWAY_URL, body) == r'mailto:test@example.com?subject=test%20subject&body=test%20body'



if __name__ == '__main__':
    unittest.main()

  

   