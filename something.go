package main

// var log logrus.FieldLogger

// func main(){
// 	log = buildTheLogger()
// 	lambda.Start(func(ctx context.Context, payload myPayload) {
// 		log = log.WithField("id", payload.ID)
// 		somethingThatUsesLog(payload)
// 	})
// }

// func main() {
// 	log = buildTheLogger()
// 	lambda.Start(func(ctx context.Context, payload myPayload) {
// 		thisInvocationLog := log.WithField("id", payload.ID)
// 		somethingThatUsesLog(thisInvocationLog, payload)
// 	})
// }

// func main() {

// 	rsp, _ := http.Get("https://www.netlify.com")
// 	pr, pw := io.Pipe()
// 	go func() {
// 		_, _ = io.Copy(pw, rsp.Body)
// 		_ = pw.Close()
// 		_ = rsp.Body.Close()
// 	}()

// 	scanner := bufio.NewScanner(pr)
// 	for scanner.Scan() {
// 		fmt.Println(scanner.Text())
// 	}
// }
