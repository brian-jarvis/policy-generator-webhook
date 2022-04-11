# Document the steps needed to deploy

1. build the app (edit image repo in Makefile if needed)
    make -f ./Makefile build

2. Log into repo
    buildah login -u="bjarvis+buildah_blog" -p="MQ6JOP4B11XEMYMPO94I8A4HK4P4ZO8ZNKNSNWRAQEVVOR3P3UO848KUZGZ5GIS5" quay.io

    buildah login -u='1979710|bjarvis-lab1' -p=eyJhbGciOiJSUzUxMiJ9.eyJzdWIiOiJiNzU3MmY2YTFjZDc0MGU0YjM2NjI0ZWYyYzM5MmJmNyJ9.LEhSs3eT-WU07mUrOgqXoRcLOP-JCWd7ef1r2G0iswMZa4dUdRsgcdjrWYHZmUiGwyYNAeiivUvWO7RuKoSSQRLB06sT8M5xBEQ2EakPCm8T6g2wa5YJaQK7wnwIWYfaR_2w5UPcLzLPKRxjDljFV7kh3c0rTks_F35npaqyTCnOIrGZLsapJmOWJkWoFPjb_G6pupsiO0D59zeGetQLTB6Y1Kn45ibsFzosW6LcuqSTxql7uNUBb-F-YipQg5irOVHN0jEq5DjzXfWSd_kDPZ8g953hPj3PnuBfw2CgYyKWnJbXC9IsiUHZVS7XaLym0INswlU_75hYJLY96PqaK5Z8r3g9etKXKEfHPMtOs68cVSZbv75PCSXngjOHeQ1c17ilAOm6N--wr2Q4pPA6DMyLig5Zi9UtH-mWvVDBIfYR8M6zscEFy3FBeX0rnp3gGw8hAo7eT0j4aG6PPTG6nhvMYPE1-Zz4R885mKAertQntQjqQCRoTDK-3IA7ED94KH1Xih0gGe0jPWOA_A9JKCBf5HlRyl6vdmwLzM8Pp_BxVZM41no9HfOpdTocPwAXO-yh5MAtZQJwjFXckbCLkMoaR_VYy5wUZvHQwWMVG4SShG4IZJPrOQ6xaKvhEHkF2Cwcflj2atudRtszjoCiUN4RhSN9u3qU6mLXU89hLk8 registry.redhat.io

3. build the image and push to repo (edit image repo in Makefile if needed)
   make -f ./Makefile push-image

## Deploy the webhook
1. Create a new project for the webhook resources
   oc new-project  mutating-webhook

2. Deploy rbac.yml
  `oc auth reconcile -f rbac.yml`

3. Deploy the webhook CRD
4. Deploy the serviceaccount
5. Deploy the serivce
6. Deploy the nginx configmap
7. Deploy the mutating configmap
8. Deploy the daemonset
9. Deploy the apiservice
10. Deploy the mutatingwebhook


## Test it
11. Create project to test in
12. apply label to namespace
    1.  webhook.toleration-mutate: enabled
13. create pods
14. verify they have the mutation

